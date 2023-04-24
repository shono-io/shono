package reaktors

import (
	"context"
	"fmt"
	"github.com/shono-io/go-shono/events"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
)

func DefaultErrorHandler(topic string, partition int32, err error) {
	if partition != -1 {
		logrus.WithField("topic", topic).WithField("partition", partition).Errorf("error while consuming: %v", err)
	} else {
		logrus.Errorf("error while consuming: %v", err)
	}
}

type ErrorHandler func(topic string, partition int32, err error)

type RunnerOpt func(*Runner)

func WithErrorHandler(eh ErrorHandler) RunnerOpt {
	return func(r *Runner) {
		r.eh = eh
	}
}

func NewRunner(rr ReaktorRegistry, er events.EventRegistry, kafka *kgo.Client, opts ...RunnerOpt) *Runner {
	r := &Runner{
		rr:    rr,
		er:    er,
		kafka: kafka,
		eh:    DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

type Runner struct {
	rr ReaktorRegistry
	er events.EventRegistry

	kafka *kgo.Client

	eh ErrorHandler
}

func (r *Runner) handleRecord(ctx context.Context, record *kgo.Record) error {
	if record == nil {
		return nil
	}

	kind := events.EventKindFromHeader(record.Headers)
	if kind == nil {
		return fmt.Errorf("missing event kind header")
	}

	event, err := r.er.Event(ctx, *kind)
	if err != nil {
		return fmt.Errorf("no event found for kind %s: %w", kind, err)
	}

	evt := event.NewInstance()
	if err := event.Decode(record.Value, &evt); err != nil {
		return fmt.Errorf("failed to decode event %s: %w", kind, err)
	}

	logrus.Debugf("received event %s", kind)

	hndlr, err := r.rr.ReaktorFor(ctx, *kind)
	if err != nil {
		return fmt.Errorf("unable to find reaktor for event kind %s: %w", kind, err)
	}

	if hndlr == nil {
		return nil
	}

	rctx := &reaktorContext{
		Context: context.Background(),
		kc:      r.kafka,
		er:      r.er,
		record:  record,
	}

	hndlr.Fn(rctx, string(record.Key), evt)

	return nil
}

func (r *Runner) Run(ctx context.Context) {
	logrus.Info("reaktor runner started")
	for {
		fetches := r.kafka.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return
		}

		if err := fetches.Err(); err != nil {
			fetches.EachError(r.eh)
			continue
		}

		for _, record := range fetches.Records() {
			err := r.handleRecord(ctx, record)
			if err != nil {
				r.eh(record.Topic, record.Partition, err)
				continue
			}

			if err := r.kafka.CommitRecords(context.Background(), record); err != nil {
				r.eh(record.Topic, record.Partition, err)
			}
		}
	}
}
