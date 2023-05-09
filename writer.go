package go_shono

import (
	"context"
	"fmt"
	"github.com/shono-io/go-shono/events"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Writer interface {
	Write(evt EventMeta, key string, payload any) error
}

type kafkaWriter struct {
	kc  *kgo.Client
	org string
}

func (w *kafkaWriter) Write(evt EventMeta, key string, payload any) error {
	val, err := evt.Encode(payload)
	if err != nil {
		return err
	}

	record := &kgo.Record{
		Key:   []byte(key),
		Value: val,
		Topic: fmt.Sprintf("%s.%s", w.org, evt.Domain),
		Headers: []kgo.RecordHeader{
			{Key: events.KindHeader, Value: []byte(evt.String())},
		},
	}

	if pr := w.kc.ProduceSync(context.Background(), record); pr.FirstErr() != nil {
		return pr.FirstErr()
	}

	return nil
}
