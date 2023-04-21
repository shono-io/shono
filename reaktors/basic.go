package reaktors

import "github.com/shono-io/go-shono/events"

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

func (b *BasicReaktorRegistry) ReaktorFor(kind events.Kind) *ReaktorInfo {
	reaktor, ok := b.reaktors[kind]
	if !ok {
		return nil
	}

	return reaktor
}
