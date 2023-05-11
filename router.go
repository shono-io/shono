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
	reaktors map[EventId]*Reaktor

	pph PoisonPillHandler
}

func NewRouter(opts ...RouterOpt) *Router {
	res := &Router{
		events:   make(map[EventId]*EventMeta),
		reaktors: make(map[EventId]*Reaktor),
		pph:      DefaultPoisonPillHandler,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func (r *Router) Register(reaktor Reaktor) {
	for _, em := range reaktor.Listen {
		if _, fnd := r.reaktors[em.EventId]; fnd {
			// -- we cannot support multiple reaktors for the same event since we cannot guarantee all
			// -- reaktors will be called.
			panic(fmt.Sprintf("reaktor for event %s already registered", em.EventId))
		}

		r.reaktors[em.EventId] = &reaktor
		r.events[em.EventId] = em
	}
}

func (r *Router) Process(ctx context.Context, eid EventId, data []byte, w Writer) {
	if data == nil {
		return
	}

	// -- get the event from the headers
	logrus.Debugf("received event %s", eid)

	// -- get the event metadata
	em, fnd := r.events[eid]
	if !fnd {
		// -- skip processing if we could not find the event
		return
	}

	// -- decode the event
	res, err := em.Decode(data)
	if err != nil {
		// !!! POISON PILL !!!
		// -- if we cannot decode the event, we cannot process it
		r.pph(ctx, eid, data, err)
		return
	}

	// -- find the reaktors for the event
	reaktor, fnd := r.reaktors[em.EventId]
	if !fnd {
		// -- skip processing if we could not find the reaktors
		return
	}

	// -- create the context
	rctx := WithEvent(ctx, em)

	reaktor.Handler(rctx, res, w)
}
