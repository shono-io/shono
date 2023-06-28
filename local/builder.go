package local

import (
	"fmt"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
)

func NewInventory() *InventoryBuilder {
	return &InventoryBuilder{
		environment: &Inventory{
			scopes:     map[string]*inventory.Scope{},
			concepts:   map[string]*inventory.Concept{},
			events:     map[string]*inventory.Event{},
			injectors:  map[string]*inventory.Injector{},
			extractors: map[string]*inventory.Extractor{},
			reactors:   map[string]*inventory.Reactor{},
		},
	}
}

type InventoryBuilder struct {
	environment *Inventory
}

func (e *InventoryBuilder) Scope(scope *inventory.Scope) *InventoryBuilder {
	e.environment.scopes[scope.Reference().String()] = scope
	return e
}

func (e *InventoryBuilder) Concept(concept *inventory.Concept) *InventoryBuilder {
	e.environment.concepts[concept.Reference().String()] = concept
	return e
}

func (e *InventoryBuilder) Event(event *inventory.Event) *InventoryBuilder {
	e.environment.events[event.Reference().String()] = event
	return e
}

func (e *InventoryBuilder) Injector(injector *inventory.Injector) *InventoryBuilder {
	e.environment.injectors[injector.Reference().String()] = injector
	return e
}

func (e *InventoryBuilder) Extractor(extractor *inventory.Extractor) *InventoryBuilder {
	e.environment.extractors[extractor.Reference().String()] = extractor
	return e
}

func (e *InventoryBuilder) Reactor(reactor *inventory.Reactor) *InventoryBuilder {
	e.environment.reactors[reactor.Reference().String()] = reactor
	return e
}

func (e *InventoryBuilder) Build() (*Inventory, error) {
	// -- make sure all references are valid
	for _, scope := range e.environment.scopes {
		if !scope.Reference().IsValid() {
			return nil, invalidReference(scope.Reference())
		}
	}

	for _, concept := range e.environment.concepts {
		if !concept.Reference().IsValid() {
			return nil, invalidReference(concept.Reference())
		}
	}

	for _, event := range e.environment.events {
		if !event.Reference().IsValid() {
			return nil, invalidReference(event.Reference())
		}
	}

	for _, injector := range e.environment.injectors {
		if !injector.Reference().IsValid() {
			return nil, invalidReference(injector.Reference())
		}
	}

	for _, extractor := range e.environment.extractors {
		if !extractor.Reference().IsValid() {
			return nil, invalidReference(extractor.Reference())
		}
	}

	for _, reactor := range e.environment.reactors {
		if !reactor.Reference().IsValid() {
			return nil, invalidReference(reactor.Reference())
		}
	}

	return e.environment, nil
}

func invalidReference(ref commons.Reference) error {
	return fmt.Errorf("invalid reference: %s", ref.String())
}
