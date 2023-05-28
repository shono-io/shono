package event

import (
	"fmt"
	"strings"
)

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

func (e EventId) Scope() string {
	return e.part(1)
}

func (e EventId) Concept() string {
	return e.part(2)
}

func (e EventId) Code() string {
	return e.part(3)
}
