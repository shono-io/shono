package local

import (
	"context"
	"github.com/shono-io/shono"
)

type scopeRepo struct {
	scopes map[string]shono.Scope
}

func (s *scopeRepo) GetScope(ctx context.Context, fqn string) (shono.Scope, bool, error) {
	res, fnd := s.scopes[fqn]
	return res, fnd, nil
}

func (s *scopeRepo) AddScope(ctx context.Context, scope shono.Scope) error {
	s.scopes[scope.FQN()] = scope
	return nil
}

func (s *scopeRepo) RemoveScope(ctx context.Context, fqn string) error {
	delete(s.scopes, fqn)
	return nil
}
