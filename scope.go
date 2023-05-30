package shono

import "github.com/shono-io/shono/logic"

type Scope interface {
	Entity
	Concept(code string, opts ...ConceptOpt) Concept
	Reaktor(code string, inputEvent EventId, logic logic.Logic, opts ...ReaktorOpt) Reaktor
}

func NewScope(code, name, description string) Scope {
	return &scope{
		NewEntity(code, code, name, description),
	}
}

type scope struct {
	Entity
}

func (s *scope) Concept(code string, opts ...ConceptOpt) Concept {
	return NewConcept(s.Code(), code, opts...)
}

func (s *scope) Reaktor(code string, inputEvent EventId, logic logic.Logic, opts ...ReaktorOpt) Reaktor {
	return NewReaktor(s.Code(), code, inputEvent, logic, opts...)
}
