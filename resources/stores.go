package resources

import "context"

type Cursor interface {
	HasNext() bool
	Next(v any) error
	Count() int64
	Close() error
}

type KeyValueStore interface {
	Get(ctx context.Context, key string, target any) error
	Set(ctx context.Context, key string, value any) error
	Delete(ctx context.Context, key string) error
}

type DocumentStore[Q any] interface {
	MustQuery(ctx context.Context, query Q) Cursor
	Query(ctx context.Context, query Q) (Cursor, error)

	MustGet(ctx context.Context, kind string, key string, target any) bool
	Get(ctx context.Context, kind string, key string, target any) (bool, error)

	MustSet(ctx context.Context, kind string, key string, value any)
	Set(ctx context.Context, kind string, key string, value any) error

	MustDelete(ctx context.Context, kind string, key ...string)
	Delete(ctx context.Context, kind string, key ...string) error
}
