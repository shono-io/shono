package local

import (
	"context"
	"github.com/shono-io/go-shono/shono"
)

func NewConceptRepo() shono.ConceptRepo {
	return &conceptRepo{
		concepts: make(map[string]shono.Concept),
	}
}

type conceptRepo struct {
	concepts map[string]shono.Concept
}

func (c *conceptRepo) GetConcept(ctx context.Context, code string) (shono.Concept, bool, error) {
	res, fnd := c.concepts[code]
	return res, fnd, nil
}

func (c *conceptRepo) AddConcept(ctx context.Context, concept shono.Concept) error {
	c.concepts[concept.GetCode()] = concept
	return nil
}

func (c *conceptRepo) RemoveConcept(ctx context.Context, code string) error {
	delete(c.concepts, code)
	return nil
}
