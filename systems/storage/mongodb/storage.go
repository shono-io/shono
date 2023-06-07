package mongodb

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
)

func NewStorage(key commons.Key, uri string, db string, opts ...graph.StorageOpt) (*graph.Storage, error) {
	return graph.NewStorage(key, "mongodb", map[string]interface{}{
		UriConf:      uri,
		DatabaseConf: db,
	}, opts...)
}
