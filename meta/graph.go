package meta

import (
	"context"
	"github.com/shono-io/shono"
)

type Client interface {
	ScopeRepo
	ReaktorRepo
	ConceptRepo
	EventRepo

	Close()
}

type ScopeRepo interface {
	GetScope(ctx context.Context, fqn string) (shono.Scope, bool, error)
	AddScope(ctx context.Context, scope shono.Scope) error
	RemoveScope(ctx context.Context, fqn string) error
}

type ReaktorRepo interface {
	GetReaktor(ctx context.Context, fqn string) (shono.Reaktor, bool, error)
	AddReaktor(ctx context.Context, reaktor shono.Reaktor) error
	RemoveReaktor(ctx context.Context, fqn string) error
}

type ConceptRepo interface {
	GetConcept(ctx context.Context, fqn string) (shono.Concept, bool, error)
	AddConcept(ctx context.Context, concept shono.Concept) error
	RemoveConcept(ctx context.Context, fqn string) error
}

type EventRepo interface {
	GetEvent(ctx context.Context, fqn string) (shono.Event, bool, error)
	AddEvent(ctx context.Context, event shono.Event) error
	RemoveEvent(ctx context.Context, fqn string) error
}
