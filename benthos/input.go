package benthos

import (
	"context"
	"github.com/shono-io/shono/graph"
)

func (g *Generator) generateInput(ctx context.Context, result map[string]any, env graph.Environment, scope graph.Scope, concept graph.Concept, reaktors []graph.Reaktor) (err error) {
	var events []graph.Key
	for _, reaktor := range reaktors {
		events = append(events, reaktor.InputEventKey())
	}

	result["input"], err = g.bb.GetConsumerConfig(g.group, events)
	return err
}
