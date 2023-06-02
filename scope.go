package shono

import "github.com/shono-io/shono/logic"

type Scope interface {
	Entity
	NewConcept(code string, opts ...ConceptOpt) Concept
	NewReaktor(code string, inputEvent EventId, logic logic.Logic, opts ...ReaktorOpt) Reaktor

	Reaktors() []Reaktor
}

func NewScope(code string, opts ...ScopeOpt) Scope {
	result := &scope{
		newEntity(NewKey("scope", code)),
		map[string]Reaktor{},
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
	reaktors map[string]Reaktor
}

func (s *scope) NewConcept(code string, opts ...ConceptOpt) Concept {
	return NewConcept(s.Key(), code, opts...)
}

func (s *scope) NewReaktor(code string, inputEvent EventId, logic logic.Logic, opts ...ReaktorOpt) Reaktor {
	result := NewReaktor(s.Key(), code, inputEvent, logic, opts...)

	s.reaktors[result.Key().String()] = result

	return result
}

func (s *scope) Reaktors() []Reaktor {
	var result []Reaktor

	for _, reaktor := range s.reaktors {
		result = append(result, reaktor)
	}

	return result
}
