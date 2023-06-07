package api

import (
	"context"
	"github.com/shono-io/shono/graph"
)

type Generator struct {
}

func (g *Generator) Generate(ctx context.Context, env graph.Environment) error {
	scope, err := env.ListScopes()
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateForScope(ctx context.Context, env graph.Environment, scope graph.Scope) error {

	return nil
}

func (g *Generator) generateForConcept(ctx context.Context, env graph.Environment, concept graph.Concept) error {

	return nil
}
