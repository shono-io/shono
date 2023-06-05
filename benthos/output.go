package benthos

import (
	"context"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
)

func (g *Generator) generateOutput(ctx context.Context, result map[string]any, env graph.Environment, scope graph.Scope, concept graph.Concept, reaktors []graph.Reaktor) (err error) {
	// -- get the list of events from the reaktors
	var events []commons.Key
	for _, reaktor := range reaktors {
		events = append(events, reaktor.OutputEventKeys()...)
	}

	bb, err := env.GetBackbone()
	if err != nil {
		return err
	}

	result["output"], err = bb.GetProducerConfig(events)
	return err
}
