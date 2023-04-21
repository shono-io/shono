package app

import (
	"context"
	"fmt"
	"github.com/shono-io/go-shono/events"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
)

func NewBackbone(cat *Catalogs, kafkaOpts ...kgo.Opt) (*Backbone, error) {
	topics := make([]string, 0, len(cat.Handlers.Handlers()))
	for _, handler := range cat.Handlers.Handlers() {
		topics = append(topics, handler.Kind().Domain)
	}
	kafkaOpts = append(kafkaOpts, kgo.ConsumeTopics(topics...))

	kafka, err := kgo.NewClient(kafkaOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %w", err)
	}

	logrus.Debugf("trying to connect to the kafka cluster")
	if err := kafka.Ping(context.Background()); err != nil { // check connectivity to cluster
		return nil, fmt.Errorf("failed to connect to kafka cluster: %w", err)
	}

	return &Backbone{
		kafka: kafka,
		cat:   cat,
	}, nil
}

type Backbone struct {
	kafka *kgo.Client
	cat   *Catalogs
}

func (h *Backbone) handleRecord(record *kgo.Record) {
	if record == nil {
		return
	}

	kind := events.EventKindFromHeader(record.Headers)
	if kind == nil {
		h.handleError(record.Topic, record.Partition, fmt.Errorf("missing event kind header"))
		return
	}

	event, fnd := h.cat.Events.Event(*kind)
	if !fnd {
		h.handleError(record.Topic, record.Partition, fmt.Errorf("no event found for event kind %s", kind))
		return
	}

	evt := event.NewInstance()
	if err := event.Decode(record.Value, &evt); err != nil {
		h.handleError(record.Topic, record.Partition, fmt.Errorf("failed to decode event %s: %w", kind, err))
		return
	}

	logrus.Debugf("received event %s", kind)

	hndlr, fnd := h.cat.Handlers.Handler(*kind)
	if !fnd {
		h.handleError(record.Topic, record.Partition, fmt.Errorf("no handler found for event kind %s", kind))
		return
	}

	ctx := &Context{
		Context:      context.Background(),
		Kind:         *kind,
		EventCatalog: h.cat.Events,
		handler:      hndlr,
		cl:           h.kafka,
	}

	if err := hndlr.Fn(ctx, string(record.Key), evt); err != nil {
		h.handleError(record.Topic, record.Partition, fmt.Errorf("failed to handle event %s: %w", kind, err))
		return
	}

}

func (h *Backbone) handleError(topic string, partition int32, err error) {
	if partition != -1 {
		logrus.WithField("topic", topic).WithField("partition", partition).Errorf("error while consuming: %v", err)
	} else {
		logrus.Errorf("error while consuming: %v", err)
	}
}

func (h *Backbone) Execute() {
	for {
		fetches := h.kafka.PollFetches(context.Background())
		if fetches.IsClientClosed() {
			return
		}

		fetches.EachError(h.handleError)
		fetches.EachRecord(h.handleRecord)

		if err := h.kafka.CommitUncommittedOffsets(context.Background()); err != nil {
			h.handleError("", -1, err)
		}
	}
}
