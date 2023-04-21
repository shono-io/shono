package events

import (
	"context"
	"errors"
)

type BasicRegistryOpt func(*BasicRegistry)

func WithEvents(events ...*EventInfo) BasicRegistryOpt {
	return func(s *BasicRegistry) {
		for _, event := range events {
			s.events[event.Kind()] = event
		}
	}
}

func WithEvent(kind Kind, opts ...EventOpt) BasicRegistryOpt {
	return func(s *BasicRegistry) {
		s.events[kind] = NewEventInfo(kind, opts...)
	}
}

func NewBasicRegistry(opts ...BasicRegistryOpt) *BasicRegistry {
	s := &BasicRegistry{
		events: make(map[Kind]*EventInfo),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type BasicRegistry struct {
	events map[Kind]*EventInfo
}

var ErrEventNotFound = errors.New("event not found")

func (s *BasicRegistry) Event(ctx context.Context, kind Kind) (*EventInfo, error) {
	event, ok := s.events[kind]
	if !ok {
		return nil, ErrEventNotFound
	}

	return event, nil
}
