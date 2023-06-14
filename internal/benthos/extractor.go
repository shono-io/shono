package benthos

import (
	"fmt"
	"github.com/shono-io/shono/core"
)

type ExtractorGenerator struct {
}

func (e *ExtractorGenerator) Generate(env core.Environment, extractorRef core.Reference) (core.Artifact, error) {
	extractor, err := env.ResolveExtractor(extractorRef)
	if err != nil {
		return nil, err
	}

	inp, err := generateBackboneInput(env, extractor.InputEvents())
	if err != nil {
		return nil, fmt.Errorf("input: %w", err)
	}

	processors, err := generateLogic(env, extractor.Logic())
	if err != nil {
		return nil, fmt.Errorf("processors: %w", err)
	}

	out, err := generateSystemOutput(env, extractor.TargetSystem(), extractor.TargetSystemConfig())
	if err != nil {
		return nil, fmt.Errorf("output: %w", err)
	}

	return NewArtifact(extractorRef.Parent(), core.ArtifactTypeExtractor, map[string]any{
		"input": inp,
		"pipeline": map[string]any{
			"processors": processors,
		},
		"output": out,
	})
}
