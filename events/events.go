package events

import (
	"context"
	"encoding/json"
	"fmt"
	sdk "github.com/shono-io/go-shono"
	"github.com/twmb/franz-go/pkg/sr"
)

func NewEvents() *Events {
	return &Events{
		events: make(map[sdk.EventKind]Event),
	}
}

type Events struct {
	events map[sdk.EventKind]Event
	src    *sr.Client
}

func (e *Events) AddEvent(event Event) error {
	s, err := e.src.Schemas(context.Background(), event.EventKind.String(), sr.HideDeleted)
	if err != nil && err != sr.ErrNotRegistered {
		return fmt.Errorf("unable to get schemas: %w", err)
	}

	schema := sr.Schema{
		Schema: string(event.Schema),
		Type:   sr.TypeJSON,
	}

	var schemaID int
	if s == nil || len(s) == 0 {
		// -- new schema, we should register it
		ss, err := e.src.CreateSchema(context.Background(), event.EventKind.String(), schema)
		if err != nil {
			return fmt.Errorf("unable to create schema: %w", err)
		}

		schemaID = ss.ID
	} else {
		ss, err := e.src.LookupSchema(context.Background(), event.EventKind.String(), schema)
		if err != nil {
			return fmt.Errorf("unable to lookup schema: %w", err)
		}

		schemaID = ss.ID
	}

	event.serde = &sr.Serde{}
	event.serde.Register(schemaID, event.valueType,
		sr.EncodeFn(func(a any) ([]byte, error) {
			return json.Marshal(a)
		}), sr.DecodeFn(func(b []byte, a any) error {
			return json.Unmarshal(b, a)
		}))

	// -- check if the schema is compatible with the existing schema in the registry
	ok, err := e.src.CheckCompatibility(context.Background(), event.EventKind.String(), -2, sr.Schema{
		Schema: string(event.Schema),
		Type:   sr.TypeJSON,
	})
	if err != nil {
		return fmt.Errorf("unable to check compatibility: %w", err)
	}

	if !ok {
		return fmt.Errorf("schema is not compatible with the existing schema")
	}

	e.events[event.EventKind] = event

	return nil
}

func (e *Events) Event(eventKind sdk.EventKind) (Event, bool) {
	event, ok := e.events[eventKind]
	return event, ok
}
