package backbone

import (
	"context"
	"fmt"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/event"
	"github.com/shono-io/go-shono/utils"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
	"strings"
	"time"
)

var (
	KafkaKeyMeta       = "kafka_key"
	KafkaTimestampMeta = "kafka_timestamp"
	KafkaTopicMeta     = "kafka_topic"
	KafkaPartitionMeta = "kafka_partition"
)

func BasedOn(evt *event.Event) event.Opt {
	return func(e *event.Event) {
		// -- copy the non-kafka metadata
		for _, k := range evt.MetaKeys() {
			if !strings.HasPrefix(k, "kafka_") {
				event.WithHeaderFromEvent(evt, k)
			}
		}

		// -- copy the message key
		event.WithHeaderFromEvent(evt, KafkaKeyMeta)(e)
	}
}

func WithKey(key string) event.Opt {
	return func(e *event.Event) {
		event.WithMetadataField(KafkaKeyMeta, key)(e)
	}
}

func NewKafkaBackbone(id string, catalog event.Catalog, opts ...kgo.Opt) Backbone {
	opts = append(opts, kgo.ConsumerGroup(id), kgo.DisableAutoCommit())

	return &KafkaBackbone{
		id:      id,
		opts:    opts,
		catalog: catalog,
		w: &Waiter{
			callbacks: map[string]callback{},
		},
		failOnEventTypeNotFound:     false,
		failOnEventIdHeaderNotFound: true,
	}
}

type KafkaBackbone struct {
	id   string
	opts []kgo.Opt

	kc *kgo.Client

	routes  []*Route
	catalog event.Catalog

	w *Waiter

	failOnEventIdHeaderNotFound bool
	failOnEventTypeNotFound     bool
}

func (b *KafkaBackbone) MustWrite(ctx context.Context, evt *event.Event) {
	if err := b.Write(ctx, evt); err != nil {
		panic(err)
	}
}

func (b *KafkaBackbone) Write(ctx context.Context, evt *event.Event) error {
	if b.kc == nil {
		return fmt.Errorf("backbone not initialized")
	}

	rec := b.eventToRecord(evt)
	rec.Topic = fmt.Sprintf("%s.%s", evt.Organization(), evt.Scope())

	if pr := b.kc.ProduceSync(ctx, rec); pr.FirstErr() != nil {
		return pr.FirstErr()
	}

	return nil
}

func (b *KafkaBackbone) WaitFor(correlationId string, timeout time.Duration, possibleEvents ...*event.EventType) (*event.Event, error) {
	return b.w.WaitFor(correlationId, timeout, possibleEvents...)
}

func (b *KafkaBackbone) Listen() error {
	if b.kc != nil {
		return fmt.Errorf("runner already running")
	}

	// -- register the event types in teh catalog
	topicMap := make(map[string]struct{})
	for _, rte := range b.routes {
		if err := b.catalog.RegisterEventType(rte.EventType); err != nil {
			return fmt.Errorf("failed to register event type %q: %w", rte.EventType.EventId, err)
		}

		topicMap[rte.EventLog] = struct{}{}

		logrus.Infof("registered event type %q", rte.EventType.EventId)
	}

	// -- add the topics we are listening to
	var topics []string
	for topic := range topicMap {
		topics = append(topics, topic)
	}
	b.opts = append(b.opts, kgo.ConsumeTopics(topics...))

	// -- create the kgo client
	kc, err := kgo.NewClient(b.opts...)
	if err != nil {
		return fmt.Errorf("failed to create kafka client: %Writer", err)
	}
	b.kc = kc

	ctx := context.Background()

	go func() {
		for {
			fetches := kc.PollFetches(ctx)
			if fetches.IsClientClosed() {
				return
			}

			// -- panic if we were unable to fetch messages
			if err := fetches.Err(); err != nil {
				panic(fmt.Errorf("failed to fetch: %w", err))
			}

			for _, record := range fetches.Records() {
				if record.Key == nil || len(record.Key) == 0 {
					logrus.Tracef("received record with no key")
					continue
				}

				logrus.Tracef("received record: %s", string(record.Key))
				err := b.handleRecord(ctx, record)
				if err != nil {
					panic(fmt.Errorf("failed to handle record: %Writer", err))
				}
			}
		}
	}()

	return nil
}

func (b *KafkaBackbone) On(eventType *event.EventType, fn EventHandler, opts ...RouteOpt) Backbone {
	rte := &Route{
		EventType: eventType,
		EventLog:  fmt.Sprintf("%s.%s", eventType.Organization(), eventType.Scope()),
		Handler:   fn,
	}

	for _, opt := range opts {
		opt(rte)
	}

	b.routes = append(b.routes, rte)

	return b
}

func (b *KafkaBackbone) Close() {
	if b.kc != nil {
		b.kc.Close()
	}
}

func (b *KafkaBackbone) handleRecord(ctx context.Context, msg *kgo.Record) error {
	if msg == nil {
		return nil
	}

	// -- convert the record into an event
	evt, err := b.recordToEvent(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to convert record to event: %w", err)
	}

	if evt == nil {
		logrus.Tracef("no event type found for %q", utils.Header(msg.Headers, go_shono.KindHeader))
		return nil
	}

	// -- notify the waiter
	b.w.handleEvent(evt)

	// -- find the route for the event
	for _, rte := range b.routes {
		if rte.EventLog != msg.Topic {
			continue
		}

		if rte.EventType.EventId != *evt.EventId {
			continue
		}

		err := rte.Handler(ctx, evt, logrus.
			WithField("topic", msg.Topic).
			WithField("event", evt.EventId).
			Logger)

		if err != nil {
			panic(fmt.Errorf("failed to handle event: %w", err))
		}

		_ = b.kc.CommitRecords(ctx, msg)
	}

	//// -- check if there is a callback for the correlation id
	//if cb, fnd := r.callbacks[corId]; fnd {
	//	if cb.matches(event.EventId(eid)) {
	//		// -- decode the message
	//		res, _, err := r.r.Decode(event.EventId(eid), msg.Record.Value)
	//		if err != nil {
	//			return fmt.Errorf("failed to resolve callback: failed to decode message: %v", err)
	//		}
	//
	//		cb.resolve(event.EventId(eid), res)
	//	}
	//}

	return nil
}

func (b *KafkaBackbone) recordToEvent(ctx context.Context, msg *kgo.Record) (*event.Event, error) {
	eventId := utils.Header(msg.Headers, go_shono.KindHeader)
	if eventId == "" {
		if b.failOnEventIdHeaderNotFound {
			return nil, fmt.Errorf("no event kind header found")
		}

		logrus.Warnf("no event kind header found")
		return nil, nil
	}

	// -- look for the event type in the catalog
	evtType, err := b.catalog.GetEventType(event.EventId(eventId))
	if err != nil {
		if b.failOnEventTypeNotFound {
			return nil, fmt.Errorf("failed to get event type %q: %w", eventId, err)
		}

		logrus.Warnf("failed to get event type %q: %v", eventId, err)
		return nil, nil
	}

	// -- skip if no eventType is found (but trace a warning)
	if evtType == nil {
		return nil, nil
	}

	// -- parse the metadata
	md := map[string]any{
		KafkaKeyMeta:       string(msg.Key),
		KafkaTimestampMeta: &msg.Timestamp,
		KafkaTopicMeta:     msg.Topic,
		KafkaPartitionMeta: msg.Partition,
	}

	// -- add all headers to the metadata as well
	for _, h := range msg.Headers {
		md[h.Key] = string(h.Value)
	}

	evt := event.NewEvent(
		evtType,
		msg.Value,
		event.WithMetadata(md),
		event.WithAcker(func() {
			if err := b.kc.CommitRecords(ctx, msg); err != nil {
				logrus.Panicf("failed to commit record with key %q: %v", string(msg.Key), err)
			}
		}),
		event.WithNacker(func(reason error, redeliver bool) {
			if redeliver {
				logrus.Warnf("Not commiting offset for key %q: %v", string(msg.Key), reason)
			} else {
				logrus.Panicf("failed to process record with key %q: %v", string(msg.Key), reason)
			}
		}),
	)

	return evt, nil
}

func (b *KafkaBackbone) eventToRecord(evt *event.Event) *kgo.Record {
	// -- create the record
	record := &kgo.Record{
		Key:   []byte(evt.MetaString(KafkaKeyMeta)),
		Value: evt.Raw(),
	}

	// -- add all metadata to the headers
	for _, k := range evt.MetaKeys() {
		if strings.HasPrefix(k, "kafka_") {
			continue
		}

		record.Headers = append(record.Headers, kgo.RecordHeader{
			Key:   k,
			Value: []byte(evt.MetaString(k)),
		})
	}

	return record
}
