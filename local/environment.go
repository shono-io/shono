package local

import (
	"github.com/shono-io/shono/backbone"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
)

func NewEnvironment(bb backbone.Backbone) *Environment {
	return &Environment{
		bb:       bb,
		scopes:   make(map[string]graph.Scope),
		concepts: make(map[string]graph.Concept),
		events:   make(map[string]graph.Event),
		reaktors: make(map[string]graph.Reaktor),
		storages: make(map[string]graph.Storage),
		stores:   make(map[string]graph.Store),
	}
}

type Environment struct {
	bb       backbone.Backbone
	scopes   map[string]graph.Scope
	concepts map[string]graph.Concept
	events   map[string]graph.Event
	reaktors map[string]graph.Reaktor
	storages map[string]graph.Storage
	stores   map[string]graph.Store
}

func (e *Environment) GetBackbone() (backbone.Backbone, error) {
	return e.bb, nil
}

func (e *Environment) GetScope(scopeKey commons.Key) (*graph.Scope, error) {
	res, fnd := e.scopes[scopeKey.String()]
	if !fnd {
		return nil, nil
	}
	return &res, nil
}

func (e *Environment) RegisterScope(scope graph.Scope) error {
	e.scopes[scope.Key().String()] = scope
	return nil
}

func (e *Environment) ListScopes() ([]graph.Scope, error) {
	var res []graph.Scope
	for _, scope := range e.scopes {
		res = append(res, scope)
	}

	return res, nil
}

func (e *Environment) GetConcept(conceptKey commons.Key) (*graph.Concept, error) {
	res, fnd := e.concepts[conceptKey.String()]
	if !fnd {
		return nil, nil
	}

	return &res, nil
}

func (e *Environment) RegisterConcept(concept graph.Concept) error {
	e.concepts[concept.Key().String()] = concept
	return nil
}

func (e *Environment) ListConceptsForScope(scopeKey commons.Key) ([]graph.Concept, error) {
	var res []graph.Concept
	for _, concept := range e.concepts {
		res = append(res, concept)
	}

	return res, nil
}

func (e *Environment) GetEvent(eventKey commons.Key) (*graph.Event, error) {
	res, fnd := e.events[eventKey.String()]
	if !fnd {
		return nil, nil
	}
	return &res, nil
}

func (e *Environment) RegisterEvent(event ...graph.Event) error {
	for _, evt := range event {
		e.events[evt.Key().String()] = evt
	}

	return nil
}

func (e *Environment) ListEventsForConcept(conceptKey commons.Key) ([]graph.Event, error) {
	var res []graph.Event
	for _, event := range e.events {
		res = append(res, event)
	}

	return res, nil
}

func (e *Environment) GetReaktor(reaktorKey commons.Key) (*graph.Reaktor, error) {
	res, fnd := e.reaktors[reaktorKey.String()]
	if !fnd {
		return nil, nil
	}
	return &res, nil
}

func (e *Environment) RegisterReaktor(reaktor ...graph.Reaktor) error {
	for _, rkt := range reaktor {
		e.reaktors[rkt.Key().String()] = rkt
	}

	return nil
}

func (e *Environment) ListReaktorsForConcept(conceptKey commons.Key) ([]graph.Reaktor, error) {
	var res []graph.Reaktor
	for _, reaktor := range e.reaktors {
		res = append(res, reaktor)
	}

	return res, nil
}

func (e *Environment) GetStorage(storageKey commons.Key) (*graph.Storage, error) {
	res, fnd := e.storages[storageKey.String()]
	if !fnd {
		return nil, nil
	}
	return &res, nil
}

func (e *Environment) RegisterStorage(storage graph.Storage) error {
	e.storages[storage.Key().String()] = storage
	return nil
}

func (e *Environment) ListStorages() ([]graph.Storage, error) {
	var res []graph.Storage
	for _, storage := range e.storages {
		res = append(res, storage)
	}

	return res, nil
}

func (e *Environment) GetStore(storeKey commons.Key) (*graph.Store, error) {
	res, fnd := e.stores[storeKey.String()]
	if !fnd {
		return nil, nil
	}
	return &res, nil
}

func (e *Environment) RegisterStore(store graph.Store) error {
	e.stores[store.Key().String()] = store
	return nil
}

func (e *Environment) ListStoresForStorage(storageKey commons.Key) ([]graph.Store, error) {
	var res []graph.Store
	for _, store := range e.stores {
		res = append(res, store)
	}

	return res, nil
}
