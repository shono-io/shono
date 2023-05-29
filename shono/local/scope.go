package local

import (
	"context"
	"github.com/shono-io/go-shono/shono"
)

func NewScopeRepo() shono.ScopeRepo {
	return &scopeRepo{
		scopes: make(map[string]shono.Scope),
	}
}

type scopeRepo struct {
	scopes map[string]shono.Scope
}

func (s *scopeRepo) GetScope(ctx context.Context, scopeCode string) (shono.Scope, bool, error) {
	res, fnd := s.scopes[scopeCode]
	return res, fnd, nil
}

func (s *scopeRepo) AddScope(ctx context.Context, scope shono.Scope) error {
	s.scopes[scope.GetCode()] = scope
	return nil
}

func (s *scopeRepo) RemoveScope(ctx context.Context, scopeCode string) error {
	delete(s.scopes, scopeCode)
	return nil
}
