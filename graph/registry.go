package graph

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
	SetScope(scope Scope) error
	SetConcept(concept Concept) error
	SetEvent(event Event) error
	SetReaktor(reaktor Reaktor) error
	On(b *ReaktorBuilder) error
}
