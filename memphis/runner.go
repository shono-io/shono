package memphis

import (
	"context"
	"fmt"
	"github.com/memphisdev/memphis.go"
	go_shono "github.com/shono-io/go-shono"
	"github.com/sirupsen/logrus"
)

func NewRunner(name string, r *go_shono.Router, host string, username string, option ...memphis.Option) (*Runner, error) {
	c, err := memphis.Connect(host, username, option...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to memphis: %w", err)
	}

	return &Runner{
		name:           name,
		c:              c,
		w:              NewWriter(name, c),
		r:              r,
		scopeConsumers: make(map[string]*memphis.Consumer),
	}, nil
}

type Runner struct {
	name string
	c    *memphis.Conn
	w    *Writer

	r              *go_shono.Router
	scopeConsumers map[string]*memphis.Consumer
}

func (r *Runner) Close() {
	for _, c := range r.scopeConsumers {
		c.StopConsume()
	}
}

func (r *Runner) createConsumer(scope string) error {
	if _, fnd := r.scopeConsumers[scope]; fnd {
		return nil
	}

	c, err := r.c.CreateConsumer(scope, r.name)
	if err != nil {
		return err
	}

	r.scopeConsumers[scope] = c
	return nil
}

func (r *Runner) handler(msgs []*memphis.Msg, err error, ctx context.Context) {
	if err != nil && err.Error() != "memphis: timeout" {
		logrus.Errorf("failed to consume: %v", err)
		return
	}

	for _, msg := range msgs {
		// -- check if there is a kind header
		eid, fnd := msg.GetHeaders()[go_shono.KindHeader]
		if !fnd {
			msg.Ack()
			logrus.Errorf("no event kind header found")
			continue
		}

		// -- get the reaktor for the event
		r.r.Process(context.Background(), go_shono.EventId(eid), msg.Data(), NewWriter("reaktor", r.c))

		// -- ack the message
		msg.Ack()
	}
}

func (r *Runner) Run() error {
	// -- create the consumers for each scope
	for _, scope := range r.r.Scopes() {
		if err := r.createConsumer(scope); err != nil {
			return fmt.Errorf("failed to create consumer: %w", err)
		}
	}

	// -- start the consumers
	for scope, c := range r.scopeConsumers {
		logrus.Infof("starting consumer for scope %s", scope)
		if err := c.Consume(r.handler); err != nil {
			return fmt.Errorf("failed to start the consumer for scope %s: %w", scope, err)
		}
	}

	return nil
}
