package graph

type Environment interface {
	GetScope(scopeKey Key) (*Scope, error)
	RegisterScope(scope Scope) error
	ListScopes() ([]Scope, error)

	GetConcept(conceptKey Key) (*Concept, error)
	RegisterConcept(concept Concept) error
	ListConceptsForScope(scopeKey Key) ([]Concept, error)

	GetEvent(eventKey Key) (*Event, error)
	RegisterEvent(event ...Event) error
	ListEventsForConcept(conceptKey Key) ([]Event, error)

	GetReaktor(reaktorKey Key) (*Reaktor, error)
	RegisterReaktor(reaktor ...Reaktor) error
	ListReaktorsForConcept(conceptKey Key) ([]Reaktor, error)

	GetStorage(storageKey Key) (*Storage, error)
	RegisterStorage(storage Storage) error
	ListStorages() ([]Storage, error)

	GetStore(storeKey Key) (*Store, error)
	RegisterStore(store Store) error
	ListStoresForStorage(storageKey Key) ([]Store, error)
}
