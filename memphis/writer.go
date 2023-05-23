package memphis

import (
	"context"
	"fmt"
	"github.com/memphisdev/memphis.go"
	go_shono "github.com/shono-io/go-shono"
)

func NewWriter(name string, c *memphis.Conn) *Writer {
	return &Writer{name, c}
}

type Writer struct {
	name string
	c    *memphis.Conn
}

func (w *Writer) MustWrite(ctx context.Context, correlationId string, evt *go_shono.EventMeta, payload go_shono.Payload) {
	if err := w.Write(ctx, correlationId, evt, payload); err != nil {
		panic(err)
	}
}

func (w *Writer) Write(ctx context.Context, correlationId string, evt *go_shono.EventMeta, payload go_shono.Payload) error {
	// -- station
	station := fmt.Sprintf("%s.%s", evt.Organization(), evt.Space())

	h := memphis.Headers{}
	h.New()
	h.Add(go_shono.KindHeader, string(evt.EventId))
	h.Add(go_shono.CorrelationHeader, correlationId)
	if err := w.c.Produce(station, w.name, payload, []memphis.ProducerOpt{}, []memphis.ProduceOpt{memphis.MsgHeaders(h)}); err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}
	return nil
}
