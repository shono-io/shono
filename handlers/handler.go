package handlers

import (
	"github.com/shono-io/go-shono/events"
	"github.com/shono-io/go-shono/states"
	"github.com/shono-io/go-shono/storage"
)

type HandlerFunc func(ctx Context, key string, value any) error

func NewHandler(kind events.Kind, fn HandlerFunc, opts ...HandlerOpt) Handler {
	h := Handler{
		kind:   kind,
		Fn:     fn,
		stores: map[states.Kind]storage.StateStore{},
	}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

type Handler struct {
	kind   events.Kind
	Fn     HandlerFunc
	emits  []events.Kind
	stores map[states.Kind]storage.StateStore
}

func (h Handler) Kind() events.Kind {
	return h.kind
}

func (h Handler) Emits() []events.Kind {
	return h.emits
}

func (h Handler) Stores() map[states.Kind]storage.StateStore {
	return h.stores
}
