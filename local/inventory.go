package local

import (
	"fmt"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
)

type Inventory struct {
	scopes     map[string]*inventory.Scope
	concepts   map[string]*inventory.Concept
	events     map[string]*inventory.Event
	injectors  map[string]*inventory.Injector
	extractors map[string]*inventory.Extractor
	reactors   map[string]*inventory.Reactor
}

func (e *Inventory) ResolveScope(ref commons.Reference) (*inventory.Scope, error) {
	res, fnd := e.scopes[ref.String()]
	if !fnd {
		return nil, fmt.Errorf("scope %s not found", ref.String())
	}
	return res, nil
}

func (e *Inventory) ResolveConcept(ref commons.Reference) (*inventory.Concept, error) {
	res, fnd := e.concepts[ref.String()]
	if !fnd {
		return nil, fmt.Errorf("concept %s not found", ref.String())
	}
	return res, nil
}

func (e *Inventory) ResolveEvent(ref commons.Reference) (*inventory.Event, error) {
	res, fnd := e.events[ref.String()]
	if !fnd {
		return nil, fmt.Errorf("event %s not found", ref.String())
	}
	return res, nil
}

func (e *Inventory) ResolveInjector(ref commons.Reference) (*inventory.Injector, error) {
	res, fnd := e.injectors[ref.String()]
	if !fnd {
		return nil, fmt.Errorf("injector %s not found", ref.String())
	}
	return res, nil
}

func (e *Inventory) ResolveExtractor(ref commons.Reference) (*inventory.Extractor, error) {
	res, fnd := e.extractors[ref.String()]
	if !fnd {
		return nil, fmt.Errorf("extractor %s not found", ref.String())
	}
	return res, nil
}

func (e *Inventory) ListReactorsForConcept(conceptRef commons.Reference) ([]inventory.Reactor, error) {
	var res []inventory.Reactor
	for _, v := range e.reactors {
		if v.Concept.String() == conceptRef.String() {
			res = append(res, *v)
		}
	}
	return res, nil
}

func (e *Inventory) ListInjectorsForScope(scopeRef commons.Reference) ([]inventory.Injector, error) {
	var res []inventory.Injector
	for _, v := range e.injectors {
		if v.Scope.String() == scopeRef.String() {
			res = append(res, *v)
		}
	}
	return res, nil
}

func (e *Inventory) ListExtractorsForScope(scopeRef commons.Reference) ([]inventory.Extractor, error) {
	var res []inventory.Extractor
	for _, v := range e.extractors {
		if v.Scope.String() == scopeRef.String() {
			res = append(res, *v)
		}
	}
	return res, nil
}

func (e *Inventory) ListScopes() ([]inventory.Scope, error) {
	var res []inventory.Scope
	for _, v := range e.scopes {
		res = append(res, *v)
	}
	return res, nil
}

func (e *Inventory) ListConceptsForScope(scopeRef commons.Reference) ([]inventory.Concept, error) {
	var res []inventory.Concept
	for _, v := range e.concepts {
		if v.Scope.String() == scopeRef.String() {
			res = append(res, *v)
		}
	}
	return res, nil
}

func (e *Inventory) ListEventsForConcept(conceptRef commons.Reference) ([]inventory.Event, error) {
	var res []inventory.Event
	for _, v := range e.events {
		if v.Concept.String() == conceptRef.String() {
			res = append(res, *v)
		}
	}
	return res, nil
}
