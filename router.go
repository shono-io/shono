package go_shono

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

type PoisonPillHandler func(ctx context.Context, eid EventId, data []byte, err error)

func DefaultPoisonPillHandler(ctx context.Context, eid EventId, data []byte, err error) {
	logrus.WithField("eid", eid).Panicf("Encountered a poison pill while consuminmg: %v", err)
}

type RouterOpt func(r *Router)

func WithPoisonPillHandler(pph PoisonPillHandler) RouterOpt {
	return func(r *Router) {
		r.pph = pph
	}
}

type Router struct {
	events   map[EventId]*EventMeta
	reaktors map[EventId][]*Reaktor

	pph PoisonPillHandler
}

func NewRouter(opts ...RouterOpt) *Router {
	res := &Router{
		events:   make(map[EventId]*EventMeta),
		reaktors: make(map[EventId][]*Reaktor),
		pph:      DefaultPoisonPillHandler,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func (r *Router) Register(reaktor Reaktor) {
	for _, em := range reaktor.Listen {
		if _, fnd := r.reaktors[em.EventId]; !fnd {
			r.reaktors[em.EventId] = make([]*Reaktor, 0)
		}

		r.reaktors[em.EventId] = append(r.reaktors[em.EventId], &reaktor)
		r.events[em.EventId] = em
	}
}

func (r *Router) Process(ctx context.Context, eid EventId, data []byte) {
	res, em, err := r.Decode(eid, data)
	if err != nil {
		// !!! POISON PILL !!!
		// -- if we cannot decode the event, we cannot process it
		r.pph(ctx, eid, data, err)
		return
	}

	if res == nil {
		return
	}

	// -- get the event from the headers
	logrus.Debugf("received event %s", eid)

	// -- find the reaktors for the event
	reaktors, fnd := r.reaktors[em.EventId]
	if !fnd {
		// -- skip processing if we could not find the reaktors
		return
	}

	// -- create the context
	rctx := WithEvent(ctx, em)

	for _, reaktor := range reaktors {
		go reaktor.Handler(rctx, res)
	}
}

func (r *Router) Decode(eid EventId, data []byte) (any, *EventMeta, error) {
	if data == nil {
		return nil, nil, nil
	}
	// -- get the event metadata
	em, fnd := r.events[eid]
	if !fnd {
		// -- skip processing if we could not find the event
		return nil, em, nil
	}

	// -- decode the event
	res, err := em.Decode(data)
	return res, em, err
}

func (r *Router) Scopes() []string {
	var scopes = make(map[string]bool)
	for _, rs := range r.reaktors {
		for _, s := range rs {
			for _, eid := range s.Listen {
				scopes[fmt.Sprintf("%s.%s", eid.Organization(), eid.Space())] = true
			}
		}
	}

	var result []string
	for scope := range scopes {
		result = append(result, scope)
	}

	return result
}
