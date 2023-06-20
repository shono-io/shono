package arangodb

import "github.com/shono-io/shono/inventory"

func NewStorage(opts ...Opt) inventory.Storage {
	config := map[string]any{}
	for _, opt := range opts {
		opt(config)
	}

	return inventory.Storage{
		Name:       "arangodb",
		ConfigSpec: configFields,
		Config:     config,
	}
}
