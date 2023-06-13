package local

import "github.com/shono-io/shono/graph"

type storageRepo struct {
	storage map[string]graph.Storage
}

func (s *storageRepo) GetStorage(key string) graph.Storage {
	return s.storage[key]
}

func (s *storageRepo) ListStorages() []graph.Storage {
	var result []graph.Storage
	for _, v := range s.storage {
		result = append(result, v)
	}
	return result
}
