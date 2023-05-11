package memphis

import (
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

func (w *Writer) MustWrite(evt *go_shono.EventMeta, payload any) {
	if err := w.Write(evt, payload); err != nil {
		panic(err)
	}
}

func (w *Writer) Write(evt *go_shono.EventMeta, payload any) error {
	// -- station
	station := fmt.Sprintf("%s.%s", evt.Organization(), evt.Space())

	h := memphis.Headers{}
	h.New()
	h.Add(go_shono.KindHeader, string(evt.EventId))
	if err := w.c.Produce(station, w.name, payload, []memphis.ProducerOpt{}, []memphis.ProduceOpt{memphis.MsgHeaders(h)}); err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}
	return nil
}
