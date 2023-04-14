package reaktor

import (
	"context"
	"fmt"
	sdk "github.com/shono-io/go-shono"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
)

func NewReaktor(handlers Handlers, config sdk.AppConfig) (*Reaktor, error) {
	opts := []kgo.Opt{
		kgo.ConsumeTopics(handlers.Topics()...),
	}
	opts = append(opts, config.KafkaOpts()...)

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("trying to connect to the kafka cluster")
	if err = cl.Ping(context.Background()); err != nil { // check connectivity to cluster
		return nil, err
	}

	return &Reaktor{
		cl:       cl,
		handlers: handlers,
	}, nil
}

type Reaktor struct {
	cl       *kgo.Client
	handlers Handlers
}

func (r *Reaktor) Start() error {
	go r.loop()

	return nil
}

func (r *Reaktor) Stop() error {
	r.cl.Close()

	return nil
}

func (r *Reaktor) handleRecord(record *kgo.Record) {
	if record == nil {
		return
	}

	kind := EventKindFromHeader(record.Headers)
	if kind == nil {
		r.handleError(record.Topic, record.Partition, fmt.Errorf("missing event kind header"))
		return
	}

	logrus.Debugf("received event %s", kind)

	h, fnd := r.handlers.Handler(*kind)
	if !fnd {
		r.handleError(record.Topic, record.Partition, fmt.Errorf("no handler found for event kind %s", kind))
		return
	}

	ctx := Context{
		Kind:        h.Kind,
		HandlerInfo: h.HandlerInfo,
	}

	if err := h.Fn(ctx, record); err != nil {
		r.handleError(record.Topic, record.Partition, fmt.Errorf("error while handling event: %s", err))
	}
}

func (r *Reaktor) handleError(topic string, partition int32, err error) {
	if partition != -1 {
		logrus.WithField("topic", topic).WithField("partition", partition).Errorf("error while consuming: %v", err)
	} else {
		logrus.Errorf("error while consuming: %v", err)
	}
}

func (r *Reaktor) loop() {
	for {
		fetches := r.cl.PollFetches(context.Background())
		if fetches.IsClientClosed() {
			return
		}

		fetches.EachError(r.handleError)
		fetches.EachRecord(r.handleRecord)

		if err := r.cl.CommitUncommittedOffsets(context.Background()); err != nil {
			r.handleError("", -1, err)
		}
	}
}
