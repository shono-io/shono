package memphis

import (
	"context"
	"fmt"
	"github.com/memphisdev/memphis.go"
	go_shono "github.com/shono-io/go-shono"
	"github.com/sirupsen/logrus"
	"time"
)

type CallbackNotification struct {
	Event    go_shono.EventId
	Value    *any
	timedOut bool
}

type callback struct {
	c      chan CallbackNotification
	events []go_shono.EventId
	t      *time.Timer
	cancel chan any
}

func (c *callback) matches(id go_shono.EventId) bool {
	for _, eid := range c.events {
		if eid == id {
			return true
		}
	}

	return false
}

func (c *callback) resolve(eid go_shono.EventId, evt any) {
	c.c <- CallbackNotification{
		Event: eid,
		Value: &evt,
	}
	c.t.Stop()
	close(c.cancel)
}

func NewRunner(name string, r *go_shono.Router, c *memphis.Conn) *Runner {
	return &Runner{
		name:           name,
		c:              c,
		w:              NewWriter(name, c),
		r:              r,
		scopeConsumers: make(map[string]*memphis.Consumer),
		callbacks:      make(map[string]callback),
	}
}

type Runner struct {
	name string
	c    *memphis.Conn
	w    *Writer

	r              *go_shono.Router
	scopeConsumers map[string]*memphis.Consumer

	callbacks map[string]callback
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
		pctx := context.Background()

		corId, fnd := msg.GetHeaders()[go_shono.CorrelationHeader]
		if !fnd {
			msg.Ack()
			logrus.Errorf("no correlation id header found")
			continue
		}

		pctx = go_shono.WithCorrelationId(pctx, corId)
		pctx = go_shono.WithEventTimestamp(pctx, time.Now())

		r.r.Process(pctx, go_shono.EventId(eid), msg.Data())

		// -- ack the message
		msg.Ack()

		// -- check if there is a callback for the correlation id
		if cb, fnd := r.callbacks[corId]; fnd {
			if cb.matches(go_shono.EventId(eid)) {
				// -- decode the message
				res, _, err := r.r.Decode(go_shono.EventId(eid), msg.Data())
				if err != nil {
					logrus.Errorf("failed to resolve callback: failed to decode message: %v", err)
					continue
				}

				cb.resolve(go_shono.EventId(eid), res)
			}
		}
	}
}

func (r *Runner) RegisterCallback(correlationId string, timeout time.Duration, events ...*go_shono.EventMeta) (chan CallbackNotification, error) {
	evts := make([]go_shono.EventId, len(events))
	for i, evt := range events {
		evts[i] = evt.EventId
	}

	// -- check if the callback already exists
	if _, fnd := r.callbacks[correlationId]; fnd {
		return nil, fmt.Errorf("callback already exists for correlation id %s", correlationId)
	}

	// -- create the channel
	ch := make(chan CallbackNotification)
	r.callbacks[correlationId] = callback{
		c:      ch,
		events: evts,
		t:      time.NewTimer(timeout),
		cancel: make(chan any),
	}

	// -- start a timer to remove the callback on timeout
	go func(cid string) {
		cb, fnd := r.callbacks[cid]
		if !fnd {
			return
		}

		select {
		case <-cb.cancel:
		case <-cb.t.C:
			cb.c <- CallbackNotification{
				timedOut: true,
			}
		}

		delete(r.callbacks, cid)
	}(correlationId)

	return ch, nil
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
