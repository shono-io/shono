package benthos

import (
	"fmt"
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
)

func NewExtractorGenerator() *ExtractorGenerator {
	return &ExtractorGenerator{}
}

type ExtractorGenerator struct {
}

func (e *ExtractorGenerator) Generate(inv inventory.Inventory, extractorRef commons.Reference) (artifacts.Artifact, error) {
	extractor, err := inv.ResolveExtractor(extractorRef)
	if err != nil {
		return nil, err
	}

	inp, err := generateBackboneInput(extractor.InputEvents())
	if err != nil {
		return nil, fmt.Errorf("input: %w", err)
	}

	dlq, err := generateBackboneDLQ()
	if err != nil {
		return nil, fmt.Errorf("dlq: %w", err)
	}

	logic, err := generateLogic(extractor.Logic())
	if err != nil {
		return nil, fmt.Errorf("logic: %w", err)
	}

	return NewArtifact(extractorRef.Parent(), commons.ArtifactTypeExtractor, *logic, inp, extractor.Output(), dlq, nil)
}
