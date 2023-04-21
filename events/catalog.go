package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/twmb/franz-go/pkg/sr"
)

func NewCatalog(schemaRegistry *sr.Client) (*Catalog, error) {
	result := &Catalog{
		schemaRegistry: schemaRegistry,
		events:         make(map[Kind]*EventInfo),
	}

	return result, nil
}

type Catalog struct {
	events         map[Kind]*EventInfo
	schemaRegistry *sr.Client
}

func (e *Catalog) RegisterEvent(kind Kind, t any) error {
	schema := jsonschema.Reflect(t)
	if schema == nil {
		return fmt.Errorf("failed to reflect schema for %T", t)
	}

	b, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	event := EventInfo{
		kind: kind,
		schema: sr.Schema{
			Schema:     string(b),
			Type:       sr.TypeJSON,
			References: nil,
		},
		serde:     sr.Serde{},
		valueType: t,
	}

	s, err := e.schemaRegistry.Schemas(context.Background(), event.kind.String(), sr.HideDeleted)
	if err != nil && err != sr.ErrNotRegistered {
		return fmt.Errorf("unable to get schemas: %w", err)
	}

	var schemaID int
	if s == nil || len(s) == 0 {
		// -- new schema, we should register it
		ss, err := e.schemaRegistry.CreateSchema(context.Background(), event.kind.String(), event.Schema())
		if err != nil {
			return fmt.Errorf("unable to create schema: %w", err)
		}

		schemaID = ss.ID
	} else {
		ss, err := e.schemaRegistry.LookupSchema(context.Background(), event.kind.String(), event.Schema())
		if err != nil {
			return fmt.Errorf("unable to lookup schema: %w", err)
		}

		schemaID = ss.ID
	}

	event.serde.Register(schemaID, event.valueType,
		sr.EncodeFn(func(a any) ([]byte, error) {
			return json.Marshal(a)
		}), sr.DecodeFn(func(b []byte, a any) error {
			return json.Unmarshal(b, a)
		}))

	// -- check if the schema is compatible with the existing schema in the registry
	ok, err := e.schemaRegistry.CheckCompatibility(context.Background(), event.kind.String(), -2, event.Schema())
	if err != nil {
		return fmt.Errorf("unable to check compatibility: %w", err)
	}

	if !ok {
		return fmt.Errorf("schema is not compatible with the existing schema")
	}

	e.events[event.kind] = &event

	return nil
}

func (e *Catalog) Event(eventKind Kind) (*EventInfo, bool) {
	event, ok := e.events[eventKind]
	return event, ok
}
