package inventory

import (
	"context"
	"github.com/shono-io/shono/commons"
)

func NewStorageReference(code string) commons.Reference {
	return commons.NewReference("storages", code)
}

type Storage struct {
	Name       string
	ConfigSpec []IOConfigSpecField
	Config     map[string]any
}

type StorageClient interface {
	List(ctx context.Context, collection string, filters map[string]any, paging *PagingOpts) (Cursor, error)
	Get(ctx context.Context, collection string, key string) (map[string]any, error)
	Set(ctx context.Context, collection string, key string, value map[string]any) error
	Add(ctx context.Context, collection string, key string, value map[string]any) error
	Delete(ctx context.Context, collection string, key string) error
	Close() error
}

type PagingOpts struct {
	Offset int64
	Size   int64
}

type Cursor interface {
	HasNext() bool
	Read() (map[string]any, error)
	Close() error
}

type StorageRepo interface {
	GetStorage(key string) Storage
	ListStorages() []Storage
}
