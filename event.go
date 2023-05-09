package go_shono

import (
	"encoding/json"
	"fmt"
)

func NewEvent(organization, domain, concept, code string, payload any) EventMeta {
	return EventMeta{
		Organization: organization,
		Domain:       domain,
		Concept:      concept,
		Code:         code,
		payload:      payload,
	}
}

type EventMeta struct {
	Organization string
	Domain       string
	Concept      string
	Code         string
	payload      any
}

func (e *EventMeta) String() string {
	return fmt.Sprintf("%s.%s.%s.%s", e.Organization, e.Domain, e.Concept, e.Code)
}

func (e *EventMeta) Encode(payload any) ([]byte, error) {
	return json.Marshal(payload)
}

func (e *EventMeta) Decode(value []byte) (any, error) {
	t := e.payload
	if err := json.Unmarshal(value, &t); err != nil {
		return nil, err
	}

	return t, nil
}
