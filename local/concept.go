package local

import "github.com/shono-io/shono/graph"

type conceptRepo struct {
	concepts map[string]graph.Concept
}

func (c *conceptRepo) GetConceptByReference(reference graph.ConceptReference) (*graph.Concept, error) {
	res, fnd := c.concepts[reference.String()]
	if !fnd {
		return nil, nil
	}

	return &res, nil
}

func (c *conceptRepo) GetConcept(scopeCode, code string) (*graph.Concept, error) {
	return c.GetConceptByReference(graph.ConceptReference{
		ScopeCode: scopeCode,
		Code:      code,
	})
}

func (c *conceptRepo) ListConceptsForScope(scopeCode string) ([]graph.Concept, error) {
	var res []graph.Concept
	for _, v := range c.concepts {
		if v.ScopeCode == scopeCode {
			res = append(res, v)
		}
	}

	return res, nil
}
