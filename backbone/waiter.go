package backbone

import (
	"fmt"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/event"
	"github.com/shono-io/go-shono/utils"
	"time"
)

type CallbackNotification struct {
	evt      *event.Event
	timedOut bool
}

type callback struct {
	c      chan CallbackNotification
	events []event.EventId
	t      *time.Timer
	cancel chan any
}

func (c *callback) matches(evt *event.Event) bool {
	for _, eid := range c.events {
		if eid == *evt.EventId {
			return true
		}
	}

	return false
}

func (c *callback) resolve(evt *event.Event) {
	c.c <- CallbackNotification{
		evt: evt,
	}
	c.t.Stop()
	close(c.cancel)
}

func NewWaiter() *Waiter {
	return &Waiter{
		callbacks: map[string]callback{},
	}
}

type Waiter struct {
	callbacks map[string]callback
}

func (w *Waiter) WaitFor(correlationId string, timeout time.Duration, possibleEvents ...*event.EventType) (*event.Event, error) {
	var pe []event.EventId
	for _, e := range possibleEvents {
		pe = append(pe, e.EventId)
	}

	c, err := w.RegisterCallback(correlationId, timeout, pe...)
	if err != nil {
		return nil, err
	}

	// -- wait for the channel to complete
	res := <-c
	if res.timedOut {
		return nil, utils.ErrTimeout
	} else {
		return res.evt, nil
	}
}

func (w *Waiter) RegisterCallback(correlationId string, timeout time.Duration, events ...event.EventId) (chan CallbackNotification, error) {
	// -- check if the callback already exists
	if _, fnd := w.callbacks[correlationId]; fnd {
		return nil, fmt.Errorf("callback already exists for correlation id %s", correlationId)
	}

	// -- create the channel
	ch := make(chan CallbackNotification)
	w.callbacks[correlationId] = callback{
		c:      ch,
		events: events,
		t:      time.NewTimer(timeout),
		cancel: make(chan any),
	}

	// -- start a timer to remove the callback on timeout
	go func(cid string) {
		cb, fnd := w.callbacks[cid]
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

		delete(w.callbacks, cid)
	}(correlationId)

	return ch, nil
}

func (w *Waiter) handleEvent(evt *event.Event) {
	corId := evt.MetaString(go_shono.CorrelationHeader)
	if corId == "" {
		return
	}

	// -- check if there is a callback for the correlation id
	if cb, fnd := w.callbacks[corId]; fnd {
		if cb.matches(evt) {
			cb.resolve(evt)
		}
	}
}
