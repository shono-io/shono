package inventory

func NewMemoryInventory() *Builder {
	return &Builder{
		inventory: newMemoryInventory(),
	}
}

type Builder struct {
	inventory *MemoryInventory
}

func (e *Builder) Scope(scope Scope) *Builder {
	e.inventory.scopes[scope.Reference().String()] = scope
	return e
}

func (e *Builder) Concept(concept Concept) *Builder {
	e.inventory.concepts[concept.Reference().String()] = concept
	return e
}

func (e *Builder) Event(event Event) *Builder {
	e.inventory.events[event.Reference().String()] = event
	return e
}

func (e *Builder) Injector(injector Injector) *Builder {
	e.inventory.injectors[injector.Reference().String()] = injector
	return e
}

func (e *Builder) Extractor(extractor Extractor) *Builder {
	e.inventory.extractors[extractor.Reference().String()] = extractor
	return e
}

func (e *Builder) Reactor(reactor Reactor) *Builder {
	e.inventory.reactors[reactor.Reference().String()] = reactor
	return e
}

func (e *Builder) Build() Inventory {
	return e.inventory
}
