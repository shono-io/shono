package graph

import "github.com/shono-io/shono/core"

type Registry interface {
	BackboneRepo
	ScopeRepo
	ConceptRepo
	EventRepo
	ReaktorRepo
	StorageRepo
}

type MutableRegistry interface {
	Registry
	SetScope(scope core.Scope) error
	SetConcept(concept core.Concept) error
	SetEvent(event Event) error
	SetReaktor(reaktor Reaktor) error
	On(b *ReaktorBuilder) error
}
