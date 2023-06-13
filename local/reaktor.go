package local

import "github.com/shono-io/shono/graph"

type reaktorRepo struct {
	reaktors map[string]graph.Reaktor
}

func (r *reaktorRepo) GetReaktorByReference(reference graph.ReaktorReference) (*graph.Reaktor, error) {
	res, fnd := r.reaktors[reference.String()]
	if !fnd {
		return nil, nil
	}

	return &res, nil
}

func (r *reaktorRepo) GetReaktor(scopeCode, conceptCode, reaktorCode string) (*graph.Reaktor, error) {
	return r.GetReaktorByReference(graph.ReaktorReference{
		ScopeCode:   scopeCode,
		ConceptCode: conceptCode,
		Code:        reaktorCode,
	})
}

func (r *reaktorRepo) ListReaktorsForConcept(scopeCode, conceptCode string) ([]graph.Reaktor, error) {
	var res []graph.Reaktor
	for _, v := range r.reaktors {
		if v.ScopeCode == scopeCode && v.ConceptCode == conceptCode {
			res = append(res, v)
		}
	}
	return res, nil
}
