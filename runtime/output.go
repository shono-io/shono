package runtime

import (
	"context"
	"github.com/shono-io/shono/graph"
)

func (g *Generator) generateOutput(ctx context.Context, result map[string]any, reg graph.Registry, scope graph.Scope, concept graph.Concept, reaktors []graph.Reaktor) (err error) {
	// -- get the list of events from the reaktors
	var events []graph.EventReference
	for _, reaktor := range reaktors {
		for _, output := range reaktor.Outputs {
			events = append(events, output.Event)
		}
	}

	bb, err := reg.GetBackbone()
	if err != nil {
		return err
	}

	res, err := bb.GetProducerConfig(events)
	if err != nil {
		return err
	}

	result["output"] = res
	return err
}
