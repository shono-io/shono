package reaktor

import (
	sdk "github.com/shono-io/go-shono"
)

type HandlerRepo interface {
	Handler(eventKind sdk.EventKind) (Handler, bool)

	TopicsToSubscribeTo() []string
}

func NewHandlers() *Handlers {
	return &Handlers{
		handlers: make(map[sdk.EventKind]Handler),
	}
}

type Handlers struct {
	handlers map[sdk.EventKind]Handler
}

func (h *Handlers) AddHandler(eventKind sdk.EventKind, info HandlerInfo, fn HandlerFunc) {
	h.handlers[eventKind] = Handler{
		HandlerInfo: info,
		Kind:        eventKind,
		Fn:          fn,
	}
}

func (h *Handlers) Handler(eventKind sdk.EventKind) (Handler, bool) {
	handler, ok := h.handlers[eventKind]
	return handler, ok
}

func (h *Handlers) Topics() []string {
	topics := make([]string, 0, len(h.handlers))

	for _, handler := range h.handlers {
		topics = append(topics, handler.Kind.Domain)
	}

	return topics
}
