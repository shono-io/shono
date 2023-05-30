package local

import (
	"context"
	"github.com/shono-io/shono"
)

type conceptRepo struct {
	concepts map[string]shono.Concept
}

func (c *conceptRepo) GetConcept(ctx context.Context, fqn string) (shono.Concept, bool, error) {
	res, fnd := c.concepts[fqn]
	return res, fnd, nil
}

func (c *conceptRepo) AddConcept(ctx context.Context, concept shono.Concept) error {
	c.concepts[concept.FQN()] = concept
	return nil
}

func (c *conceptRepo) RemoveConcept(ctx context.Context, fqn string) error {
	delete(c.concepts, fqn)
	return nil
}
