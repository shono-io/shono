package go_shono

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/twmb/franz-go/pkg/sr"
	"strings"
)

func NewEvent(eid EventId, payload any) *EventMeta {
	return &EventMeta{
		EventId: eid,
		payload: payload,
	}
}

func NewEventId(organization, space, concept, code string) EventId {
	return EventId(fmt.Sprintf("%s:%s:%s:%s", organization, space, concept, code))
}

type EventId string

func (e EventId) part(i int) string {
	parts := strings.Split(string(e), ":")
	if len(parts) != 4 {
		return ""
	}
	return parts[i]
}

func (e EventId) Organization() string {
	return e.part(0)
}

func (e EventId) Space() string {
	return e.part(1)
}

func (e EventId) Concept() string {
	return e.part(2)
}

func (e EventId) Code() string {
	return e.part(3)
}

type EventMeta struct {
	EventId
	payload any
	serde   sr.Serde
}

func (e *EventMeta) Register(client *sr.Client) error {
	// -- get the schema of the event
	s := jsonschema.Reflect(e.payload)
	sb, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("unable to marshal schema: %w", err)
	}

	sch := sr.Schema{Type: sr.TypeJSON, Schema: string(sb)}

	// -- check if we can create the schema
	ss, err := client.CreateSchema(context.Background(), string(e.EventId), sch)
	if err != nil {
		return fmt.Errorf("unable to register schema: %w", err)
	}

	e.serde.Register(ss.ID, e.payload, sr.EncodeFn(json.Marshal), sr.DecodeFn(json.Unmarshal))

	return nil
}

func (e *EventMeta) Encode(payload any) ([]byte, error) {
	return e.serde.Encode(payload)
}

func (e *EventMeta) Decode(value []byte) (any, error) {
	t := e.payload
	if err := e.serde.Decode(value, &t); err != nil {
		return nil, err
	}

	return t, nil
}
