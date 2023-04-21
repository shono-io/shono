package reaktors

import (
	"context"
	"github.com/shono-io/go-shono/events"
)

type ReaktorRegistry interface {
	ReaktorFor(ctx context.Context, kind events.Kind) (*ReaktorInfo, error)
	Topics(ctx context.Context) []string
}
