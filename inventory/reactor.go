package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewReactorReference(scopeCode, conceptCode, code string) commons.Reference {
	return NewConceptReference(scopeCode, conceptCode).Child("reactors", code)
}

type Reactor struct {
	Node
	Concept          commons.Reference
	InputEvent       commons.Reference
	OutputEventCodes []string
	Logic            Logic
}

func (r *Reactor) Reference() commons.Reference {
	return NewReactorReference(r.Concept.Parent().Code(), r.Concept.Code(), r.Code)
}
