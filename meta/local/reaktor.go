package local

import (
	"context"
	"github.com/shono-io/shono"
)

type reaktorRepo struct {
	reaktors map[string]shono.Reaktor
}

func (r *reaktorRepo) GetReaktor(ctx context.Context, fqn string) (shono.Reaktor, bool, error) {
	res, fnd := r.reaktors[fqn]
	return res, fnd, nil
}

func (r *reaktorRepo) AddReaktor(ctx context.Context, reaktor shono.Reaktor) error {
	r.reaktors[reaktor.FQN()] = reaktor
	return nil
}

func (r *reaktorRepo) RemoveReaktor(ctx context.Context, fqn string) error {
	delete(r.reaktors, fqn)
	return nil
}
