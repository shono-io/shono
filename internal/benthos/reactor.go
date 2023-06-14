package benthos

import (
	"fmt"
	"github.com/shono-io/shono/core"
)

type ReaktorGenerator struct {
}

func (g *ReaktorGenerator) Generate(env core.Environment, conceptRef core.Reference) (core.Artifact, error) {
	reaktors, err := env.ListReactorsForConcept(conceptRef)
	if err != nil {
		return nil, err
	}

	content, err := generateReactors(env, reaktors)
	if err != nil {
		return nil, err
	}

	return NewArtifact(conceptRef, core.ArtifactTypeReaktor, content)
}

func generateReactors(env core.Environment, reactors []core.Reactor) (map[string]any, error) {
	var inputEvents []core.Reference
	var outputEvents []core.Reference
	for _, r := range reactors {
		inputEvents = append(inputEvents, r.InputEvent())
		outputEvents = append(outputEvents, r.OutputEvents()...)
	}

	inp, err := generateBackboneInput(env, inputEvents)
	if err != nil {
		return nil, fmt.Errorf("input: %w", err)
	}

	processors, err := generateProcessors(env, reactors)
	if err != nil {
		return nil, fmt.Errorf("processors: %w", err)
	}

	out, err := generateBackboneOutput(env, outputEvents)
	if err != nil {
		return nil, fmt.Errorf("output: %w", err)
	}

	return map[string]any{
		"input": inp,
		"pipeline": map[string]any{
			"processors": []map[string]any{
				{
					"switch": processors,
				},
			},
		},
		"output": out,
	}, nil
}

func generateProcessors(env core.Environment, reactors []core.Reactor) ([]map[string]any, error) {
	var cases []map[string]any

	// -- for each reaktor, we will add a case checking the type header
	for _, reactor := range reactors {
		processors, err := generateLogic(env, reactor.Logic())
		if err != nil {
			return nil, fmt.Errorf("logic for reactor %q failed to generate: %w", reactor.Reference(), err)
		}

		if processors != nil && len(processors) > 0 {
			cases = append(cases, map[string]any{
				"check":      fmt.Sprintf("@io_shono_kind == %q", reactor.InputEvent()),
				"processors": processors,
			})
		}
	}

	// -- we will add a trace to indicate that no processor was found
	cases = append(cases, map[string]any{
		"processors": []map[string]any{
			{
				"log": map[string]any{
					"level":   "TRACE",
					"message": "no processor for ${!this.format_json()}",
				},
			},
		},
	})

	return cases, nil
}
