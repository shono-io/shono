package storage

import (
	"context"
	"github.com/shono-io/go-shono/states"
)

type PagingParams struct {
	Size   int `json:"size"`
	Offset int `json:"offset"`
}

type ListResult struct {
	Total int64 `json:"total"`
	Items []any `json:"items"`
}

type DocumentHandler func(doc map[string]interface{}) error

type Storage interface {
	ForState(ctx context.Context, kind states.Kind) (StateStore, error)
}

type Filter interface {
}
