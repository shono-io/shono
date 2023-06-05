package graph

import "github.com/shono-io/shono/commons"

type StorageOpt func(*Storage)

func WithStorageName(name string) StorageOpt {
	return func(s *Storage) {
		s.name = name
	}
}

func WithStorageDescription(description string) StorageOpt {
	return func(s *Storage) {
		s.description = description
	}
}

func NewStorage(key commons.Key, kind string, config map[string]any, opts ...StorageOpt) (*Storage, error) {
	result := &Storage{
		key:    key,
		name:   key.Code(),
		kind:   kind,
		config: config,
	}

	for _, opt := range opts {
		opt(result)
	}

	return result, nil
}

type Storage struct {
	key         commons.Key
	name        string
	description string
	kind        string
	config      map[string]any
}

func (s Storage) Key() commons.Key {
	return s.key
}

func (s Storage) Name() string {
	return s.name
}

func (s Storage) Description() string {
	return s.description
}

func (s Storage) Kind() string {
	return s.kind
}

func (s Storage) Config() map[string]any {
	return s.config
}
