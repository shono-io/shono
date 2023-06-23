package local

import (
	"github.com/shono-io/shono/inventory"
)

func NewInventory() *InventoryBuilder {
	return &InventoryBuilder{
		environment: &Inventory{
			scopes:     map[string]inventory.Scope{},
			concepts:   map[string]inventory.Concept{},
			events:     map[string]inventory.Event{},
			injectors:  map[string]inventory.Injector{},
			extractors: map[string]inventory.Extractor{},
			reactors:   map[string]inventory.Reactor{},
		},
	}
}

type InventoryBuilder struct {
	environment *Inventory
}

func (e *InventoryBuilder) Scope(scope inventory.Scope) *InventoryBuilder {
	e.environment.scopes[scope.Reference().String()] = scope
	return e
}

func (e *InventoryBuilder) Concept(concept inventory.Concept) *InventoryBuilder {
	e.environment.concepts[concept.Reference().String()] = concept
	return e
}

func (e *InventoryBuilder) Event(event inventory.Event) *InventoryBuilder {
	e.environment.events[event.Reference().String()] = event
	return e
}

func (e *InventoryBuilder) Injector(injector inventory.Injector) *InventoryBuilder {
	e.environment.injectors[injector.Reference().String()] = injector
	return e
}

func (e *InventoryBuilder) Extractor(extractor inventory.Extractor) *InventoryBuilder {
	e.environment.extractors[extractor.Reference().String()] = extractor
	return e
}

func (e *InventoryBuilder) Reactor(reactor inventory.Reactor) *InventoryBuilder {
	e.environment.reactors[reactor.Reference().String()] = reactor
	return e
}

func (e *InventoryBuilder) Build() *Inventory {
	return e.environment
}
