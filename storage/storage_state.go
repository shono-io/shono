package storage

import (
	"context"
	"github.com/shono-io/go-shono/states"
)

type StateStore interface {
	Kind() states.Kind
	Exists(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string, target interface{}) (bool, error)
	Store(ctx context.Context, key string, data any) (string, error)
	Delete(ctx context.Context, key string) error
}

type DocumentStateStore interface {
	StateStore
	Query(ctx context.Context, query string, vars map[string]interface{}, handler DocumentHandler) (int64, error)
	List(ctx context.Context, filters map[string]interface{}, params *PagingParams, handler DocumentHandler) (int64, error)
	Merge(ctx context.Context, key string, data any) error
}
