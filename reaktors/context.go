package reaktors

import (
	"context"
	"fmt"
	"github.com/shono-io/go-shono/events"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

type SendOpt func(r *kgo.Record)

func WithKey(key string) SendOpt {
	return func(r *kgo.Record) {
		r.Key = []byte(key)
	}
}

type ReaktorContext interface {
	context.Context
	Key() string
	Timestamp() time.Time
	Failed(kind events.Kind, err error)
	Send(kind events.Kind, value any, opts ...SendOpt)
	Header(key string) []byte
}

type reaktorContext struct {
	context.Context
	kc     *kgo.Client
	er     events.EventRegistry
	record *kgo.Record
}

func (ctx *reaktorContext) Key() string {
	return string(ctx.record.Key)
}

func (ctx *reaktorContext) Timestamp() time.Time {
	return ctx.record.Timestamp
}

func (ctx *reaktorContext) Header(key string) []byte {
	for _, h := range ctx.record.Headers {
		if h.Key == key {
			return h.Value
		}
	}

	return nil
}

func (ctx *reaktorContext) Failed(kind events.Kind, err error) {
	ctx.Send(kind, events.OperationFailed{Reason: err.Error()})
}

func (ctx *reaktorContext) Send(kind events.Kind, value any, opts ...SendOpt) {
	ei, err := ctx.er.Event(ctx, kind)
	if err != nil {
		panic(fmt.Sprintf("error getting event info: %v", err))
	}

	val, err := ei.Encode(value)
	if err != nil {
		panic(fmt.Sprintf("error marshaling value: %v", err))
	}

	record := &kgo.Record{
		Key:   ctx.record.Key,
		Value: val,
		Topic: kind.Domain,
		Headers: []kgo.RecordHeader{
			{Key: events.KindHeader, Value: []byte(kind.String())},
		},
	}

	for _, opt := range opts {
		opt(record)
	}

	if pr := ctx.kc.ProduceSync(ctx, record); pr.FirstErr() != nil {
		panic(fmt.Sprintf("error producing record: %v", pr.FirstErr()))
	}
}
