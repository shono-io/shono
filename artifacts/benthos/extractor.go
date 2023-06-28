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

func (e *ExtractorGenerator) Generate(applicationId string, artifactId string, inv inventory.Inventory, extractorRef commons.Reference) (*artifacts.Artifact, error) {
	if applicationId == "" {
		return nil, fmt.Errorf("applicationId must be provided")
	}

	if artifactId == "" {
		return nil, fmt.Errorf("artifactId must be provided")
	}

	if !extractorRef.IsValid() {
		return nil, fmt.Errorf("invalid extractor reference")
	}

	// -- get the extractor based on the given reference
	extractor, err := inv.ResolveExtractor(extractorRef)
	if err != nil {
		return nil, fmt.Errorf("extractor: %w", err)
	}

	opts := []artifacts.Opt{
		artifacts.WithOutput(artifacts.AsGeneratedOutput(extractor.Output)),
	}

	// -- resolve the input events
	var inputEvents []artifacts.GeneratedEvent
	for _, evtRef := range extractor.InputEvents {
		evt, err := inv.ResolveEvent(evtRef)
		if err != nil {
			return nil, fmt.Errorf("input event: %w", err)
		}
		if evt == nil {
			return nil, fmt.Errorf("input event: %s not found", evtRef)
		}
		inputEvents = append(inputEvents, artifacts.AsGeneratedEvent(*evt))
	}
	opts = append(opts, artifacts.WithInputEvents(inputEvents...))

	inp, err := generateBackboneInput(applicationId, extractor.InputEvents)
	if err != nil {
		return nil, fmt.Errorf("input: %w", err)
	}
	opts = append(opts, artifacts.WithInput(*inp))

	dlq, err := generateBackboneDLQ()
	if err != nil {
		return nil, fmt.Errorf("dlq: %w", err)
	}
	opts = append(opts, artifacts.WithDlq(*dlq))

	logic, err := generateLogic(extractor.Logic)
	if err != nil {
		return nil, fmt.Errorf("logic: %w", err)
	}
	opts = append(opts, artifacts.WithLogic(*logic))

	return artifacts.NewArtifact(commons.ArtifactTypeExtractor, extractor.Reference().Parent(), fmt.Sprintf("%s_%s", applicationId, artifactId), opts...)
}
