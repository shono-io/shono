package graph

import "github.com/shono-io/shono/commons"

type StoreOpt func(s *Store)

func WithStoreName(name string) StoreOpt {
	return func(s *Store) {
		s.name = name
	}
}

func WithStoreDescription(description string) StoreOpt {
	return func(s *Store) {
		s.description = description
	}
}

func NewStore(key commons.Key, storageKey commons.Key, col string, opts ...StoreOpt) Store {
	res := Store{
		key:        key,
		storageKey: storageKey,
		name:       key.Code(),
		collection: col,
	}

	for _, opt := range opts {
		opt(&res)
	}

	return res
}

type Store struct {
	key         commons.Key
	storageKey  commons.Key
	name        string
	description string
	collection  string
}

func (s Store) Key() commons.Key {
	return s.key
}

func (s Store) Name() string {
	return s.name
}

func (s Store) Description() string {
	return s.description
}

func (s Store) Collection() string {
	return s.collection
}

func (s Store) StorageKey() commons.Key {
	return s.storageKey
}
