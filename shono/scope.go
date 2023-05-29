package shono

import "context"

// == ENTITY ==========================================================================================================

type Scope interface {
	Entity
	ConceptRepo
	ReaktorRepo
}

func NewScope(code, name, description string, conceptRepo ConceptRepo, reaktorRepo ReaktorRepo) Scope {
	return &scope{
		NewEntity(code, code, name, description),
		conceptRepo,
		reaktorRepo,
	}
}

type scope struct {
	Entity
	ConceptRepo
	ReaktorRepo
}

// == REPO ============================================================================================================

type ScopeRepo interface {
	GetScope(ctx context.Context, scopeCode string) (Scope, bool, error)
	AddScope(ctx context.Context, scope Scope) error
	RemoveScope(ctx context.Context, scopeCode string) error
}
