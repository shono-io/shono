package event

import go_shono "github.com/shono-io/go-shono"

func NewEventType(eid EventId, payload any, serde go_shono.Serde) *EventType {
	return &EventType{
		EventId: eid,
		payload: payload,
		serde:   serde,
	}
}

type EventType struct {
	EventId
	payload any
	serde   go_shono.Serde
}

//func (e *EventType) Register(client *sr.Client) error {
//	// -- get the schema of the event
//	s := jsonschema.Reflect(e.payload)
//	sb, err := json.Marshal(s)
//	if err != nil {
//		return fmt.Errorf("unable to marshal schema: %w", err)
//	}
//
//	sch := sr.Schema{Type: sr.TypeJSON, Schema: string(sb)}
//
//	// -- check if we can create the schema
//	ss, err := client.CreateSchema(context.Background(), string(e.EventId), sch)
//	if err != nil {
//		return fmt.Errorf("unable to register schema: %w", err)
//	}
//
//	e.serde.Register(ss.ID, e.payload, sr.EncodeFn(json.Marshal), sr.DecodeFn(json.Unmarshal))
//
//	return nil
//}

func (e *EventType) Encode(payload any) ([]byte, error) {
	return e.serde.Encode(payload)
}

func (e *EventType) Decode(value []byte, target any) error {
	if err := e.serde.Decode(value, &target); err != nil {
		return err
	}

	return nil
}
