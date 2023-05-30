package shono

import (
	"fmt"
)

type Event interface {
	Entity
	Id() EventId
}

type EventOpt func(*event)

func WithEventName(name string) EventOpt {
	return func(e *event) {
		e.entity.name = name
	}
}

func WithEventDescription(description string) EventOpt {
	return func(e *event) {
		e.entity.description = description
	}
}

func NewEvent(scopeCode, conceptCode, code string, opts ...EventOpt) Event {
	result := &event{
		newEntity(fmt.Sprintf("%s:%s:%s", scopeCode, conceptCode, code), code),
		scopeCode,
		conceptCode,
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type event struct {
	*entity
	scopeCode   string
	conceptCode string
}

func (e *event) Id() EventId {
	if e == nil {
		return ""
	}

	return EventId(e.FQN())
}
