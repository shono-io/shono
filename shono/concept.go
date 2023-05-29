package shono

import (
	"context"
	"fmt"
)

// == ENTITY ==========================================================================================================

type Concept interface {
	Entity
	EventRepo

	FQN() string
}

func NewConcept(scopeCode, code, name, description string, eventRepo EventRepo) Concept {
	return &concept{
		scopeCode,
		NewEntity(fmt.Sprintf("%s:%s", scopeCode, code), code, name, description),
		eventRepo,
	}
}

type concept struct {
	scopeCode string
	Entity
	EventRepo
}

// == REPO ============================================================================================================

type ConceptRepo interface {
	GetConcept(ctx context.Context, code string) (Concept, bool, error)
	AddConcept(ctx context.Context, concept Concept) error
	RemoveConcept(ctx context.Context, code string) error
}
