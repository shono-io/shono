package artifacts

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
)

type ConceptGenerator interface {
	// Generate create a new reaktor artifact for all reactors for a given concept.
	Generate(env inventory.Inventory, conceptRef commons.Reference) (Artifact, error)
}

type InjectorGenerator interface {
	// Generate create a new injector artifact for the given injector.
	Generate(env inventory.Inventory, injectorRef commons.Reference) (Artifact, error)
}

type ExtractorGenerator interface {
	// Generate create a new extractor artifact for the given extractor.
	Generate(env inventory.Inventory, extractorRef commons.Reference) (Artifact, error)
}
