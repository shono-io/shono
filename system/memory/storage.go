package memory

import "github.com/shono-io/shono/inventory"

func NewStorage(storageId string) inventory.Storage {
	return inventory.Storage{
		Name:       "memory",
		ConfigSpec: []inventory.IOConfigSpecField{},
		Config: map[string]any{
			StorageIdField: storageId,
		},
	}
}
