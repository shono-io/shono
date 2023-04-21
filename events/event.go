package events

import (
	"github.com/shono-io/go-shono/codec"
)

type EventSchema string

type EventOpt func(*EventInfo)

func WithValueType(v any) EventOpt {
	return func(e *EventInfo) {
		e.valueType = v
	}
}

func NewEventInfo(kind Kind, opts ...EventOpt) *EventInfo {
	event := &EventInfo{
		kind:      kind,
		codec:     &codec.Json{},
		valueType: map[string]any{},
	}

	for _, opt := range opts {
		opt(event)
	}

	return event
}

type EventInfo struct {
	kind      Kind
	codec     codec.Codec
	valueType any
}

func (e *EventInfo) Kind() Kind {
	return e.kind
}

func (e *EventInfo) Encode(value any) ([]byte, error) {
	return e.codec.Encode(value)
}

func (e *EventInfo) Decode(b []byte, value any) error {
	return e.codec.Decode(b, value)
}

func (e *EventInfo) NewInstance() any {
	return e.valueType
}
