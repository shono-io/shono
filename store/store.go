package store

import "context"

type PersistMode string

var (
	CreatePersistMode  PersistMode = "create"
	ReplacePersistMode PersistMode = "replace"
	PatchPersistMode   PersistMode = "patch"
)

type State interface {
	Key() string
}

type Store[T State] interface {
	List(ctx context.Context, filters map[string]interface{}, offset uint, size uint) ([]T, int64, error)
	Get(ctx context.Context, key string) (*T, error)
	Remove(ctx context.Context, key string) error
	Persist(ctx context.Context, state T, mode PersistMode) error
}
