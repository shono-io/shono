package local

import (
	"github.com/shono-io/shono/core"
	"github.com/shono-io/shono/graph"
)

type conceptRepo struct {
	concepts map[string]core.Concept
}

func (c *conceptRepo) GetConceptByReference(reference graph.ConceptReference) (*core.Concept, error) {
	res, fnd := c.concepts[reference.String()]
	if !fnd {
		return nil, nil
	}

	return &res, nil
}

func (c *conceptRepo) GetConcept(scopeCode, code string) (*core.Concept, error) {
	return c.GetConceptByReference(graph.ConceptReference{
		ScopeCode: scopeCode,
		Code:      code,
	})
}

func (c *conceptRepo) ListConceptsForScope(scopeCode string) ([]core.Concept, error) {
	var res []core.Concept
	for _, v := range c.concepts {
		if v.ScopeCode == scopeCode {
			res = append(res, v)
		}
	}

	return res, nil
}
