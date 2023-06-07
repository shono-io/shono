package systems

import "github.com/shono-io/shono/graph"

var Storage = map[string]graph.StorageSystem{}

func RegisterStorageSystem(kind string, system graph.StorageSystem) {
	Storage[kind] = system
}
