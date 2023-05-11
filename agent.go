package go_shono

import (
	"context"
	"fmt"
	"github.com/shono-io/go-shono/events"
	"github.com/shono-io/go-shono/utils"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

type AgentOpt func(a *Agent)

func WithKafkaOpts(opts ...kgo.Opt) AgentOpt {
	return func(a *Agent) {
		a.kafkaOpts = opts
	}
}

func WithSchemaRegistryOpts(opts ...sr.Opt) AgentOpt {
	return func(a *Agent) {
		a.srOpts = opts
	}
}

func WithReaktor(r ...Reaktor) AgentOpt {
	return func(a *Agent) {
		for _, reaktor := range r {
			for _, em := range reaktor.Listen {
				if _, fnd := a.reaktors[em.EventId]; fnd {
					// -- we cannot support multiple reaktors for the same event since we cannot guarantee all
					// -- reaktors will be called.
					panic(fmt.Sprintf("reaktor for event %s already registered", em.EventId))
				}

				a.reaktors[em.EventId] = &reaktor
				a.events[em.EventId] = em
			}
		}
	}
}

func WithErrorHandler(eh ErrorHandler) AgentOpt {
	return func(a *Agent) {
		a.eh = eh
	}
}

func WithPoisonPillHandler(pph PoisonPillHandler) AgentOpt {
	return func(a *Agent) {
		a.pph = pph
	}
}

func NewAgent(orgId string, appId string, opts ...AgentOpt) *Agent {
	a := &Agent{
		organization:  orgId,
		applicationId: appId,

		events:     make(map[EventId]*EventMeta),
		reaktors:   make(map[EventId]*Reaktor),
		extraktors: make(map[string]Extraktor),

		eh:  DefaultErrorHandler,
		pph: DefaultPoisonPillHandler,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

type ErrorHandler func(topic string, partition int32, err error)

func DefaultErrorHandler(topic string, partition int32, err error) {
	if partition != -1 {
		logrus.WithField("topic", topic).WithField("partition", partition).Errorf("error while consuming: %v", err)
	} else {
		logrus.Errorf("error while consuming: %v", err)
	}
}

type PoisonPillHandler func(record *kgo.Record, err error)

func DefaultPoisonPillHandler(record *kgo.Record, err error) {
	logrus.WithField("topic", record.Topic).WithField("partition", record.Partition).Panicf("Encountered a poison pill while consuminmg: %v", err)
}

type Agent struct {
	organization  string
	applicationId string
	kafkaOpts     []kgo.Opt
	srOpts        []sr.Opt

	events     map[EventId]*EventMeta
	reaktors   map[EventId]*Reaktor
	extraktors map[string]Extraktor

	eh  ErrorHandler
	pph PoisonPillHandler
}

func (a *Agent) Run() error {
	src, err := sr.NewClient(a.srOpts...)
	if err != nil {
		return fmt.Errorf("unable to create schema registry client: %w", err)
	}

	if err := a.validateSchemas(src); err != nil {
		return fmt.Errorf("unable to validate schemas: %w", err)
	}

	// -- get the topics from the reaktors
	topics := a.topics()
	a.kafkaOpts = append(a.kafkaOpts, kgo.ConsumeTopics(topics...))
	a.kafkaOpts = append(a.kafkaOpts, kgo.ConsumerGroup(a.applicationId))
	a.kafkaOpts = append(a.kafkaOpts, kgo.DisableAutoCommit())

	// -- create the client
	kc, err := kgo.NewClient(a.kafkaOpts...)
	if err != nil {
		return fmt.Errorf("unable to create kafka client: %w", err)
	}
	defer kc.Close()

	logrus.Info("agent started")
	logrus.Infof("listening to topics: %v", topics)
	ctx := context.Background()
	for {
		fetches := kc.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return nil
		}

		if err := fetches.Err(); err != nil {
			fetches.EachError(a.eh)
			continue
		}

		a.processRecords(ctx, kc, fetches.Records())
	}
}

func (a *Agent) validateSchemas(src *sr.Client) error {
	// -- for each event, register the schema
	for _, event := range a.events {
		if err := event.Register(src); err != nil {
			return fmt.Errorf("unable to register schema for event %s: %w", event.EventId, err)
		}
	}

	return nil
}

func (a *Agent) processRecords(ctx context.Context, kc *kgo.Client, records []*kgo.Record) {
	//defer func() {
	//	if err := recover(); err != nil {
	//		a.eh("", -1, fmt.Errorf("panic while processing record: %v", err))
	//	}
	//}()

	ctx = WithOrganization(ctx, a.organization)
	w := &kafkaWriter{kc: kc, org: a.organization}

	for _, record := range records {
		a.handleRecord(ctx, record, w)

		if err := kc.CommitRecords(context.Background(), record); err != nil {
			a.eh(record.Topic, record.Partition, err)
		}
	}
}

func (a *Agent) handleRecord(ctx context.Context, record *kgo.Record, w Writer) {
	if record == nil {
		return
	}

	// -- get the event from the headers
	em := a.eventFromHeader(record.Headers)
	if em == nil {
		// -- skip processing if we could not find the event
		return
	}

	logrus.Debugf("received event %s", em.EventId)

	// -- decode the event
	res, err := em.Decode(record.Value)
	if err != nil {
		// !!! POISON PILL !!!
		// -- if we cannot decode the event, we cannot process it
		a.pph(record, err)
		return
	}

	// -- find the reaktors for the event
	reaktor, fnd := a.reaktors[em.EventId]
	if !fnd {
		// -- skip processing if we could not find the reaktors
		return
	}

	// -- create the context
	rctx := WithOrganization(ctx, a.organization)
	rctx = WithKey(rctx, string(record.Key))
	rctx = WithEvent(rctx, em)

	reaktor.Handler(rctx, res, w)
}

func (a *Agent) eventFromHeader(headers []kgo.RecordHeader) *EventMeta {
	value := utils.Header(headers, events.KindHeader)
	if value == "" {
		return nil
	}

	event, fnd := a.events[EventId(value)]
	if !fnd {
		return nil
	}

	return event
}

func (a *Agent) topics() []string {
	topics := map[string]struct{}{}
	for evt, _ := range a.reaktors {
		topic := fmt.Sprintf("%s.%s", a.organization, evt.Space())
		topics[topic] = struct{}{}
	}

	var result []string
	for topic := range topics {
		result = append(result, topic)
	}

	return result
}
