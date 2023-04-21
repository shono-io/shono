package reaktors

import (
	"context"
	"fmt"
	"github.com/shono-io/go-shono/events"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

type ReaktorContext interface {
	context.Context
	Timestamp() time.Time
	Send(kind events.Kind, key string, value any) error
	Header(key string) []byte
}

type reaktorContext struct {
	context.Context
	kc     *kgo.Client
	er     events.EventRegistry
	record *kgo.Record
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

func (ctx *reaktorContext) Send(kind events.Kind, key string, value any) error {
	ei, err := ctx.er.Event(ctx, kind)
	if err != nil {
		return fmt.Errorf("error getting event info: %v", err)
	}

	val, err := ei.Encode(value)
	if err != nil {
		return fmt.Errorf("error marshaling value: %v", err)
	}

	record := &kgo.Record{
		Key:   []byte(key),
		Value: val,
		Topic: kind.Domain,
		Headers: []kgo.RecordHeader{
			{Key: "shono.type", Value: []byte(kind.String())},
		},
	}

	if pr := ctx.kc.ProduceSync(ctx, record); pr.FirstErr() != nil {
		return fmt.Errorf("error producing record: %v", pr.FirstErr())
	}

	return nil
}
