package events

import (
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
	sdk "github.com/shono-io/go-shono"
	"github.com/twmb/franz-go/pkg/sr"
)

type EventSchema string

func NewEventFromStruct(kind sdk.EventKind, t any) (*Event, error) {
	schema := jsonschema.Reflect(t)
	if schema == nil {
		return nil, fmt.Errorf("failed to reflect schema for %T", t)
	}

	b, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	return &Event{
		EventKind: kind,
		Schema:    EventSchema(b),
		valueType: t,
	}, nil
}

func NewEvent(kind sdk.EventKind, schema EventSchema) *Event {
	return &Event{
		EventKind: kind,
		Schema:    schema,
		valueType: map[string]interface{}{},
	}
}

type Event struct {
	sdk.EventKind
	Schema    EventSchema
	serde     *sr.Serde
	valueType any
}

func (e *Event) Encode(value any) []byte {
	return e.serde.MustEncode(value)
}

func (e *Event) Decode(b []byte, value any) error {
	return e.serde.Decode(b, value)
}
