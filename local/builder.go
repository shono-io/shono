package local

import (
	"github.com/shono-io/shono/inventory"
)

func NewInventory() *EnvironmentBuilder {
	return &EnvironmentBuilder{
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

type EnvironmentBuilder struct {
	environment *Inventory
}

func (e *EnvironmentBuilder) Scope(scope inventory.Scope) *EnvironmentBuilder {
	e.environment.scopes[scope.Reference().String()] = scope
	return e
}

func (e *EnvironmentBuilder) Concept(concept inventory.Concept) *EnvironmentBuilder {
	e.environment.concepts[concept.Reference().String()] = concept
	return e
}

func (e *EnvironmentBuilder) Event(event inventory.Event) *EnvironmentBuilder {
	e.environment.events[event.Reference().String()] = event
	return e
}

func (e *EnvironmentBuilder) Injector(injector inventory.Injector) *EnvironmentBuilder {
	e.environment.injectors[injector.Reference().String()] = injector
	return e
}

func (e *EnvironmentBuilder) Extractor(extractor inventory.Extractor) *EnvironmentBuilder {
	e.environment.extractors[extractor.Reference().String()] = extractor
	return e
}

func (e *EnvironmentBuilder) Reactor(reactor inventory.Reactor) *EnvironmentBuilder {
	e.environment.reactors[reactor.Reference().String()] = reactor
	return e
}

func (e *EnvironmentBuilder) Build() *Inventory {
	return e.environment
}
