package runtime

import (
	"context"
	"github.com/shono-io/shono/graph"
)

func (g *Generator) generateInput(ctx context.Context, result map[string]any, reg graph.Registry, scope graph.Scope, concept graph.Concept, reaktors []graph.Reaktor) (err error) {
	var events []graph.EventReference
	for _, reaktor := range reaktors {
		events = append(events, *reaktor.Input)
	}

	bb, err := reg.GetBackbone()
	if err != nil {
		return err
	}

	res, err := bb.GetConsumerConfig(g.group, events)
	if err != nil {
		return err
	}

	result["input"] = res

	return err
}
