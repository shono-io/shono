package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewEventReference(scopeCode, conceptCode, code string) commons.Reference {
	return NewConceptReference(scopeCode, conceptCode).Child("events", code)
}

type Event struct {
	Node
	Concept commons.Reference
}

func (e *Event) Reference() commons.Reference {
	return NewEventReference(e.Concept.Parent().Parent().Code(), e.Concept.Parent().Code(), e.Code)
}
