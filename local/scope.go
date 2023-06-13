package local

import (
	"github.com/shono-io/shono/graph"
)

type scopeRepo struct {
	scopes map[string]graph.Scope
}

func (s *scopeRepo) GetScope(code string) (*graph.Scope, error) {
	res, fnd := s.scopes[code]
	if !fnd {
		return nil, nil
	}

	return &res, nil
}

func (s *scopeRepo) ListScopes() ([]graph.Scope, error) {
	var res []graph.Scope
	for _, v := range s.scopes {
		res = append(res, v)
	}
	return res, nil
}
