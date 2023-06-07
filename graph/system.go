package graph

import (
	"context"
)

type System[T any] interface {
	GetClient(config map[string]any) (T, error)
}

type StorageSystem System[StorageClient]

type PagingOpts struct {
	Offset int64
	Size   int64
}

type Cursor interface {
	HasNext() bool
	Read() (map[string]any, error)
	Close() error
}

type StorageClient interface {
	List(ctx context.Context, collection string, filters map[string]any, paging *PagingOpts) (Cursor, error)
	Get(ctx context.Context, collection string, key string) (map[string]any, error)
	Set(ctx context.Context, collection string, key string, value map[string]any) error
	Add(ctx context.Context, collection string, key string, value map[string]any) error
	Delete(ctx context.Context, collection string, key string) error
	Close() error
}
