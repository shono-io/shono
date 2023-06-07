package benthos

import (
	"context"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
)

func (g *Generator) generateInput(ctx context.Context, result map[string]any, env graph.Environment, scope graph.Scope, concept graph.Concept, reaktors []graph.Reaktor) (err error) {
	var events []commons.Key
	for _, reaktor := range reaktors {
		events = append(events, reaktor.InputEventKey())
	}

	bb, err := env.GetBackbone()
	if err != nil {
		return err
	}

	res, err := bb.GetConsumerConfig(g.group, events)
	if err != nil {
		return err
	}

	//delete(res, "")

	result["input"] = res

	return err
}
