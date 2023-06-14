package local

import (
	"github.com/shono-io/shono/core"
)

type scopeRepo struct {
	scopes map[string]core.Scope
}

func (s *scopeRepo) GetScope(code string) (*core.Scope, error) {
	res, fnd := s.scopes[code]
	if !fnd {
		return nil, nil
	}

	return &res, nil
}

func (s *scopeRepo) ListScopes() ([]core.Scope, error) {
	var res []core.Scope
	for _, v := range s.scopes {
		res = append(res, v)
	}
	return res, nil
}
