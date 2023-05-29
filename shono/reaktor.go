package shono

import "context"

// == ENTITY ==========================================================================================================

type Reaktor interface {
	Entity
	Run(ctx context.Context) (err error)
	Close() error
}

// == REPO ============================================================================================================

type ReaktorRepo interface {
	GetReaktor(ctx context.Context, code string) (Reaktor, bool, error)
	AddReaktor(ctx context.Context, reaktor Reaktor) error
	RemoveReaktor(ctx context.Context, code string) error
}
