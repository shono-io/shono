package graph

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/systems/backbone"
)

type Environment interface {
	GetBackbone() (backbone.Backbone, error)

	GetScope(scopeKey commons.Key) (*Scope, error)
	RegisterScope(scope Scope) error
	ListScopes() ([]Scope, error)

	GetConcept(conceptKey commons.Key) (*Concept, error)
	RegisterConcept(concept Concept) error
	ListConceptsForScope(scopeKey commons.Key) ([]Concept, error)

	GetEvent(eventKey commons.Key) (*Event, error)
	RegisterEvent(event ...Event) error
	ListEventsForConcept(conceptKey commons.Key) ([]Event, error)

	GetReaktor(reaktorKey commons.Key) (*Reaktor, error)
	RegisterReaktor(reaktor ...Reaktor) error
	ListReaktorsForConcept(conceptKey commons.Key) ([]Reaktor, error)

	GetStorage(storageKey commons.Key) (*Storage, error)
	RegisterStorage(storage Storage) error
	ListStorages() ([]Storage, error)

	GetStore(storeKey commons.Key) (*Store, error)
	RegisterStore(store Store) error
	ListStoresForStorage(storageKey commons.Key) ([]Store, error)
}
