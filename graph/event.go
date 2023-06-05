package graph

import "github.com/shono-io/shono/commons"

type EventOpt func(*Event)

func WithEventName(name string) EventOpt {
	return func(e *Event) {
		e.name = name
	}
}

func WithEventDescription(description string) EventOpt {
	return func(e *Event) {
		e.description = description
	}
}

func NewEvent(key commons.Key, opts ...EventOpt) Event {
	result := Event{
		key:  key,
		name: key.Code(),
	}

	for _, opt := range opts {
		opt(&result)
	}

	return result
}

type Event struct {
	key         commons.Key
	name        string
	description string
}

func (e Event) Key() commons.Key {
	return e.key
}

func (e Event) Name() string {
	return e.name
}

func (e Event) Description() string {
	return e.description
}
