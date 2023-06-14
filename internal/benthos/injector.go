package benthos

import (
	"fmt"
	"github.com/shono-io/shono/core"
)

type InjectorGenerator struct {
}

func (i *InjectorGenerator) Generate(env core.Environment, injectorRef core.Reference) (core.Artifact, error) {
	injector, err := env.ResolveInjector(injectorRef)
	if err != nil {
		return nil, err
	}

	inp, err := generateSystemInput(env, injector.SourceSystem(), injector.SourceSystemConfig())
	if err != nil {
		return nil, fmt.Errorf("input: %w", err)
	}

	processors, err := generateLogic(env, injector.Logic())
	if err != nil {
		return nil, fmt.Errorf("processors: %w", err)
	}

	out, err := generateBackboneOutput(env, injector.OutputEvents())
	if err != nil {
		return nil, fmt.Errorf("output: %w", err)
	}

	return NewArtifact(injectorRef.Parent(), core.ArtifactTypeInjector, map[string]any{
		"input": inp,
		"pipeline": map[string]any{
			"processors": processors,
		},
		"output": out,
	})
}
