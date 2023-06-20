package inventory

import (
	"errors"
	"github.com/shono-io/shono/commons"
)

var ErrNotFound = errors.New("not found")

func newMemoryInventory() *MemoryInventory {
	return &MemoryInventory{
		scopes:     make(map[string]Scope),
		concepts:   make(map[string]Concept),
		events:     make(map[string]Event),
		injectors:  make(map[string]Injector),
		extractors: make(map[string]Extractor),
		reactors:   make(map[string]Reactor),
	}
}

type MemoryInventory struct {
	scopes     map[string]Scope
	concepts   map[string]Concept
	events     map[string]Event
	injectors  map[string]Injector
	extractors map[string]Extractor
	reactors   map[string]Reactor
}

func (e *MemoryInventory) ResolveScope(ref commons.Reference) (Scope, error) {
	res, fnd := e.scopes[ref.String()]
	if !fnd {
		return nil, ErrNotFound
	}
	return res, nil
}

func (e *MemoryInventory) ResolveConcept(ref commons.Reference) (Concept, error) {
	res, fnd := e.concepts[ref.String()]
	if !fnd {
		return nil, ErrNotFound
	}
	return res, nil
}

func (e *MemoryInventory) ResolveEvent(ref commons.Reference) (Event, error) {
	res, fnd := e.events[ref.String()]
	if !fnd {
		return nil, ErrNotFound
	}
	return res, nil
}

func (e *MemoryInventory) ResolveInjector(ref commons.Reference) (Injector, error) {
	res, fnd := e.injectors[ref.String()]
	if !fnd {
		return nil, ErrNotFound
	}
	return res, nil
}

func (e *MemoryInventory) ResolveExtractor(ref commons.Reference) (Extractor, error) {
	res, fnd := e.extractors[ref.String()]
	if !fnd {
		return nil, ErrNotFound
	}
	return res, nil
}

func (e *MemoryInventory) ListReactorsForConcept(conceptRef commons.Reference) ([]Reactor, error) {
	var res []Reactor
	for _, v := range e.reactors {
		if v.Concept().String() == conceptRef.String() {
			res = append(res, v)
		}
	}
	return res, nil
}
