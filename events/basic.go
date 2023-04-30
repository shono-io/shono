package events

import (
	"context"
	"errors"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
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

func (s *BasicRegistry) MustAsRecord(ctx context.Context, kind Kind, key string, value any) *kgo.Record {
	record, err := s.AsRecord(ctx, kind, key, value)
	if err != nil {
		panic(err)
	}

	return record
}

func (s *BasicRegistry) AsRecord(ctx context.Context, kind Kind, key string, value any) (*kgo.Record, error) {
	evt, err := s.Event(ctx, kind)
	if err != nil {
		panic(fmt.Sprintf("unknown event kind %s: %v", kind, err))
	}

	val, err := evt.Encode(value)
	if err != nil {
		panic(fmt.Sprintf("error encoding event %s: %v", kind, err))
	}

	return &kgo.Record{
		Key:   []byte(key),
		Value: val,
		Topic: kind.Domain,
		Headers: []kgo.RecordHeader{
			{Key: "io.shono.kind", Value: []byte(kind.String())},
		},
	}, nil
}

var ErrEventNotFound = errors.New("event not found")

func (s *BasicRegistry) Event(ctx context.Context, kind Kind) (*EventInfo, error) {
	event, ok := s.events[kind]
	if !ok {
		return nil, ErrEventNotFound
	}

	return event, nil
}
