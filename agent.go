package go_shono

import (
	"context"
	"fmt"
	"github.com/shono-io/go-shono/events"
	"github.com/shono-io/go-shono/utils"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
)

type AgentOpt func(a *Agent)

func WithKafkaOpts(opts ...kgo.Opt) AgentOpt {
	return func(a *Agent) {
		a.kafkaOpts = opts
	}
}

func WithReaktor(r ...Reaktor) AgentOpt {
	return func(a *Agent) {
		for _, reaktor := range r {
			for _, event := range reaktor.Listen {
				if _, fnd := a.reaktors[event]; fnd {
					// -- we cannot support multiple reaktors for the same event since we cannot guarantee all
					// -- reaktors will be called.
					panic(fmt.Sprintf("reaktor for event %s already registered", event))
				}

				a.reaktors[event] = reaktor
			}
		}
	}
}

func WithResource(res ...Resource[any]) AgentOpt {
	return func(a *Agent) {
		for _, resource := range res {
			a.runtime.resources[resource.Id] = resource
		}
	}
}

func WithErrorHandler(eh ErrorHandler) AgentOpt {
	return func(a *Agent) {
		a.eh = eh
	}
}

func NewAgent(orgId string, appId string, opts ...AgentOpt) *Agent {
	a := &Agent{
		organization:  orgId,
		applicationId: appId,

		runtime: &Runtime{
			resources: make(map[string]Resource[any]),
		},

		reaktors:   make(map[EventMeta]Reaktor),
		extraktors: make(map[string]Extraktor),

		eh: DefaultErrorHandler,
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

type Agent struct {
	organization  string
	applicationId string
	kafkaOpts     []kgo.Opt

	runtime *Runtime

	events     map[string]EventMeta
	reaktors   map[EventMeta]Reaktor
	extraktors map[string]Extraktor

	eh ErrorHandler
}

func (a *Agent) Run() error {
	// -- get the topics from the reaktors
	a.kafkaOpts = append(a.kafkaOpts, kgo.ConsumeTopics(a.topics()...))
	a.kafkaOpts = append(a.kafkaOpts, kgo.ConsumerGroup(a.applicationId))

	// -- create the client
	kc, err := kgo.NewClient(a.kafkaOpts...)
	if err != nil {
		return fmt.Errorf("unable to create kafka client: %w", err)
	}
	defer kc.Close()

	// -- set the kafka client in the runtime
	a.runtime.kc = kc

	logrus.Info("agent started")
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

func (a *Agent) processRecords(ctx context.Context, kc *kgo.Client, records []*kgo.Record) {
	defer func() {
		if err := recover(); err != nil {
			a.eh("", -1, fmt.Errorf("panic while processing record: %v", err))
		}
	}()

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
	evt := a.eventFromHeader(record.Headers)
	if evt == nil {
		// -- skip processing if we could not find the event
		return
	}

	logrus.Debugf("received event %s", evt)

	// -- find the reaktors for the event
	reaktor, fnd := a.reaktors[*evt]
	if !fnd {
		// -- skip processing if we could not find the reaktors
		return
	}

	// -- create the context
	rctx := WithOrganization(ctx, a.organization)
	rctx = WithKey(rctx, string(record.Key))
	rctx = WithEvent(rctx, *evt)

	reaktor.Handler(rctx, a.runtime, w)
}

func (a *Agent) eventFromHeader(headers []kgo.RecordHeader) *EventMeta {
	value := utils.Header(headers, events.KindHeader)
	if value == "" {
		return nil
	}

	event, fnd := a.events[value]
	if !fnd {
		return nil
	}

	return &event
}

func (a *Agent) topics() []string {
	var topics map[string]struct{}
	for evt, _ := range a.reaktors {
		topic := fmt.Sprintf("%s.%s", a.organization, evt.Domain)
		topics[topic] = struct{}{}
	}

	var result []string
	for topic := range topics {
		result = append(result, topic)
	}

	return result
}
