package local

import (
	"context"
	"github.com/shono-io/go-shono/shono"
)

func NewEventRepo() shono.EventRepo {
	return &eventRepo{
		events: make(map[string]shono.Event),
	}
}

type eventRepo struct {
	events map[string]shono.Event
}

func (e *eventRepo) GetEvent(ctx context.Context, code string) (shono.Event, bool, error) {
	res, fnd := e.events[code]
	return res, fnd, nil
}

func (e *eventRepo) AddEvent(ctx context.Context, event shono.Event) error {
	e.events[event.GetCode()] = event
	return nil
}

func (e *eventRepo) RemoveEvent(ctx context.Context, code string) error {
	delete(e.events, code)
	return nil
}
