package mongodb

import "github.com/shono-io/shono/inventory"

func NewStorage(opt ...Opt) inventory.Storage {
	config := map[string]any{}
	for _, o := range opt {
		o(config)
	}

	return inventory.Storage{
		Name:       "mongodb",
		ConfigSpec: configFields,
		Config:     config,
	}
}
