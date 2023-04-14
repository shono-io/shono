package reaktor

import (
	"context"
	"fmt"
	sdk "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/events"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

type Context struct {
	context.Context
	Kind        sdk.EventKind
	HandlerInfo HandlerInfo
	cl          *kgo.Client
	Events      *events.Events
}

func (c *Context) Send(kind sdk.EventKind, key string, value any) error {
	e, fnd := c.Events.Event(kind)
	if !fnd {
		return fmt.Errorf("event %s not found", kind)
	}

	record := &kgo.Record{
		Key:   []byte(key),
		Value: e.Encode(value),
		Headers: []kgo.RecordHeader{
			{
				Key:   sdk.KindHeader,
				Value: []byte(kind.String()),
			},
		},
		Timestamp: time.Now(),
		Topic:     kind.Domain,
	}

	pr := c.cl.ProduceSync(c.Context, record)
	if pr.FirstErr() != nil {
		return pr.FirstErr()
	}

	return nil
}
