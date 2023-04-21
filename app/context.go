package app

import (
	"context"
	"fmt"
	"github.com/shono-io/go-shono/events"
	"github.com/shono-io/go-shono/handlers"
	"github.com/shono-io/go-shono/states"
	"github.com/shono-io/go-shono/storage"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

type Context struct {
	context.Context
	Kind         events.Kind
	EventCatalog *events.Catalog
	handler      *handlers.Handler
	cl           *kgo.Client
	record       *kgo.Record
}

func (c *Context) Headers() []kgo.RecordHeader {
	return c.record.Headers
}

func (c *Context) Timestamp() time.Time {
	return c.record.Timestamp
}

func (c *Context) Topic() string {
	return c.record.Topic
}

func (c *Context) Partition() int32 {
	return c.record.Partition
}

func (c *Context) Store(kind states.Kind) (storage.StateStore, error) {
	s, fnd := c.handler.Stores()[kind]
	if !fnd {
		return nil, fmt.Errorf("state store %s not found", kind)
	}

	return s, nil
}

func (c *Context) Send(kind events.Kind, key string, value any) error {
	e, fnd := c.EventCatalog.Event(kind)
	if !fnd {
		return fmt.Errorf("event %s not found", kind)
	}

	record := &kgo.Record{
		Key:   []byte(key),
		Value: e.Encode(value),
		Headers: []kgo.RecordHeader{
			{
				Key:   events.KindHeader,
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
