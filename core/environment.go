package core

type Environment interface {
	ApplicationId() string
	Backbone() Backbone

	ResolveScope(ref Reference) (Scope, error)
	ResolveConcept(ref Reference) (Concept, error)
	ResolveEvent(ref Reference) (Event, error)
	ResolveSystem(ref Reference) (System, error)
	ResolveInjector(ref Reference) (Injector, error)
	ResolveExtractor(ref Reference) (Extractor, error)

	ListReactorsForConcept(conceptRef Reference) ([]Reactor, error)
}
