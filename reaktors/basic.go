package reaktors

import (
	"context"
	"github.com/shono-io/go-shono/events"
)

type BasicReaktorRegistryOpt func(*BasicReaktorRegistry)

func NewBasicReaktorRegistry(opts ...BasicReaktorRegistryOpt) *BasicReaktorRegistry {
	registry := &BasicReaktorRegistry{
		reaktors: map[events.Kind]*ReaktorInfo{},
	}

	for _, opt := range opts {
		opt(registry)
	}

	return registry
}

func WithReaktors(reaktors ...*ReaktorInfo) BasicReaktorRegistryOpt {
	return func(b *BasicReaktorRegistry) {
		for _, reaktor := range reaktors {
			b.reaktors[reaktor.Consumes] = reaktor
		}
	}
}

type BasicReaktorRegistry struct {
	reaktors map[events.Kind]*ReaktorInfo
}

func (b *BasicReaktorRegistry) ReaktorFor(ctx context.Context, kind events.Kind) (*ReaktorInfo, error) {
	reaktor, ok := b.reaktors[kind]
	if !ok {
		return nil, nil
	}

	return reaktor, nil
}

func (b *BasicReaktorRegistry) Topics(ctx context.Context) []string {
	topics := make([]string, 0)
	for _, reaktor := range b.reaktors {
		topics = append(topics, reaktor.Consumes.Domain)
	}

	return topics
}
