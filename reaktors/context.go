package reaktors

import (
	"context"
	"fmt"
	"github.com/shono-io/go-shono/events"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ReaktorContext struct {
	context.Context
	kc *kgo.Client
	er events.EventRegistry
}

func (h *ReaktorContext) Send(kind events.Kind, key string, value any) error {
	ei, err := h.er.Event(h, kind)
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

	if pr := h.kc.ProduceSync(h, record); pr.FirstErr() != nil {
		return fmt.Errorf("error producing record: %v", pr.FirstErr())
	}

	return nil
}
