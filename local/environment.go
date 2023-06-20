package local

import (
	"errors"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/runtime"
)

var ErrNotFound = errors.New("not found")

type Inventory struct {
	scopes     map[string]inventory.Scope
	concepts   map[string]inventory.Concept
	events     map[string]inventory.Event
	systems    map[string]runtime.System
	injectors  map[string]inventory.Injector
	extractors map[string]inventory.Extractor
	reactors   map[string]inventory.Reactor
}

func (e *Inventory) ResolveScope(ref commons.Reference) (inventory.Scope, error) {
	res, fnd := e.scopes[ref.String()]
	if !fnd {
		return nil, ErrNotFound
	}
	return res, nil
}

func (e *Inventory) ResolveConcept(ref commons.Reference) (inventory.Concept, error) {
	res, fnd := e.concepts[ref.String()]
	if !fnd {
		return nil, ErrNotFound
	}
	return res, nil
}

func (e *Inventory) ResolveEvent(ref commons.Reference) (inventory.Event, error) {
	res, fnd := e.events[ref.String()]
	if !fnd {
		return nil, ErrNotFound
	}
	return res, nil
}

func (e *Inventory) ResolveSystem(ref commons.Reference) (runtime.System, error) {
	res, fnd := e.systems[ref.String()]
	if !fnd {
		return nil, ErrNotFound
	}
	return res, nil
}

func (e *Inventory) ResolveInjector(ref commons.Reference) (inventory.Injector, error) {
	res, fnd := e.injectors[ref.String()]
	if !fnd {
		return nil, ErrNotFound
	}
	return res, nil
}

func (e *Inventory) ResolveExtractor(ref commons.Reference) (inventory.Extractor, error) {
	res, fnd := e.extractors[ref.String()]
	if !fnd {
		return nil, ErrNotFound
	}
	return res, nil
}

func (e *Inventory) ListReactorsForConcept(conceptRef commons.Reference) ([]inventory.Reactor, error) {
	var res []inventory.Reactor
	for _, v := range e.reactors {
		if v.Concept().String() == conceptRef.String() {
			res = append(res, v)
		}
	}
	return res, nil
}
