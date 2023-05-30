package local

import (
	"context"
	"github.com/shono-io/shono"
)

type eventRepo struct {
	events map[string]shono.Event
}

func (e *eventRepo) GetEvent(ctx context.Context, fqn string) (shono.Event, bool, error) {
	res, fnd := e.events[fqn]
	return res, fnd, nil
}

func (e *eventRepo) AddEvent(ctx context.Context, event shono.Event) error {
	e.events[event.FQN()] = event
	return nil
}

func (e *eventRepo) RemoveEvent(ctx context.Context, fqn string) error {
	delete(e.events, fqn)
	return nil
}
