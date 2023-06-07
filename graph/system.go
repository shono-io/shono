package graph

import (
	"context"
	"github.com/shono-io/shono/commons"
)

type System[T any] interface {
	GetClient(config map[string]any) (T, error)
}

type StorageSystem System[StorageClient]

type PagingOpts struct {
	Offset int
	Size   int
}

type Cursor interface {
	HasNext() bool
	Read() (map[string]any, error)
	Close() error
	Count() int64
}

type StorageClient interface {
	List(ctx context.Context, collection string, filters map[string]any, paging *PagingOpts) (Cursor, error)
	Get(ctx context.Context, collection string, key commons.Key) (map[string]any, error)
	Set(ctx context.Context, collection string, key commons.Key, value map[string]any) error
	Add(ctx context.Context, collection string, key commons.Key, value map[string]any) error
	Delete(ctx context.Context, collection string, key commons.Key) error
}
