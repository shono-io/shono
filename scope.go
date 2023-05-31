package shono

import "github.com/shono-io/shono/logic"

type Scope interface {
	Entity
	NewConcept(code string, opts ...ConceptOpt) Concept
	NewReaktor(code string, inputEvent EventId, logic logic.Logic, opts ...ReaktorOpt) Reaktor
}

func NewScope(code string, opts ...ScopeOpt) Scope {
	result := &scope{
		newEntity(NewKey("scope", code)),
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type ScopeOpt func(s *scope)

func WithScopeName(name string) ScopeOpt {
	return func(s *scope) {
		s.name = name
	}
}

func WithScopeDescription(description string) ScopeOpt {
	return func(s *scope) {
		s.description = description
	}
}

type scope struct {
	*entity
}

func (s *scope) NewConcept(code string, opts ...ConceptOpt) Concept {
	return NewConcept(s.Key(), code, opts...)
}

func (s *scope) NewReaktor(code string, inputEvent EventId, logic logic.Logic, opts ...ReaktorOpt) Reaktor {
	return NewReaktor(s.Key(), code, inputEvent, logic, opts...)
}
