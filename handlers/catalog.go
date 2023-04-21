package handlers

import (
	"github.com/shono-io/go-shono/events"
)

func NewCatalog() (*Catalog, error) {
	result := &Catalog{
		handlers: make(map[events.Kind]*Handler),
	}

	return result, nil
}

type Catalog struct {
	handlers map[events.Kind]*Handler
}

func (e *Catalog) RegisterHandler(handler Handler) {
	e.handlers[handler.Kind()] = &handler
}

func (e *Catalog) Handler(eventKind events.Kind) (*Handler, bool) {
	hn, fnd := e.handlers[eventKind]
	return hn, fnd
}

func (e *Catalog) Handlers() map[events.Kind]*Handler {
	return e.handlers
}
