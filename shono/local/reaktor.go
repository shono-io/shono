package local

import (
	"context"
	"github.com/shono-io/go-shono/shono"
)

func NewReaktorRepo() shono.ReaktorRepo {
	return &reaktorRepo{
		reaktors: make(map[string]shono.Reaktor),
	}
}

type reaktorRepo struct {
	reaktors map[string]shono.Reaktor
}

func (r *reaktorRepo) GetReaktor(ctx context.Context, code string) (shono.Reaktor, bool, error) {
	res, fnd := r.reaktors[code]
	return res, fnd, nil
}

func (r *reaktorRepo) AddReaktor(ctx context.Context, reaktor shono.Reaktor) error {
	r.reaktors[reaktor.GetCode()] = reaktor
	return nil
}

func (r *reaktorRepo) RemoveReaktor(ctx context.Context, code string) error {
	delete(r.reaktors, code)
	return nil
}
