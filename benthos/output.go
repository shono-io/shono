package benthos

import (
	"context"
	"github.com/shono-io/shono/graph"
)

func (g *Generator) generateOutput(ctx context.Context, result map[string]any, env graph.Environment, scope graph.Scope, concept graph.Concept, reaktors []graph.Reaktor) (err error) {
	// -- get the list of events from the reaktors
	var events []graph.Key
	for _, reaktor := range reaktors {
		events = append(events, reaktor.OutputEventKeys()...)
	}

	result["output"], err = g.bb.GetProducerConfig(events)
	return err
}
