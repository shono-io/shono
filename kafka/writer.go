package kafka

import (
	"context"
	"fmt"
	go_shono "github.com/shono-io/go-shono"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Writer struct {
	kc *kgo.Client
}

func (w *Writer) MustWrite(ctx context.Context, correlationId string, evt *go_shono.EventMeta, payload go_shono.Payload) {
	if err := w.Write(ctx, correlationId, evt, payload); err != nil {
		panic(err)
	}
}

func (w *Writer) Write(ctx context.Context, correlationId string, evt *go_shono.EventMeta, payload go_shono.Payload) error {
	val, err := evt.Encode(payload)
	if err != nil {
		return err
	}

	record := &kgo.Record{
		Key:   []byte(payload.Key()),
		Value: val,
		Topic: fmt.Sprintf("%s.%s", evt.Organization(), evt.Space()),
		Headers: []kgo.RecordHeader{
			{Key: go_shono.KindHeader, Value: []byte(evt.EventId)},
			{Key: go_shono.CorrelationHeader, Value: []byte(correlationId)},
		},
	}

	if pr := w.kc.ProduceSync(ctx, record); pr.FirstErr() != nil {
		return pr.FirstErr()
	}

	return nil
}
