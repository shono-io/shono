package handlers

import (
	"github.com/shono-io/go-shono/events"
	"github.com/shono-io/go-shono/storage"
)

type HandlerOpt func(h *Handler)

func WithEmitting(kind events.Kind) HandlerOpt {
	return func(h *Handler) {
		h.emits = append(h.emits, kind)
	}
}

func WithStorage(store storage.StateStore) HandlerOpt {
	return func(h *Handler) {
		h.stores[store.Kind()] = store
	}
}
