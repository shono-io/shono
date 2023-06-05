package arangodb

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
)

func NewStorage(key commons.Key, urls []string, username, password, database string, opts ...graph.StorageOpt) (*graph.Storage, error) {
	return graph.NewStorage(key, "arangodb", map[string]interface{}{
		"urls":     urls,
		"username": username,
		"password": password,
		"database": database,
	}, opts...)
}
