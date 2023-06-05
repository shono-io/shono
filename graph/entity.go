package graph

import "github.com/shono-io/shono/commons"

type Entity interface {
	Key() commons.Key
	Name() string
	Description() string
}

func NewEntity(key commons.Key, name, description string) Entity {
	return &entity{
		key:         key,
		name:        name,
		description: description,
	}
}

func newEntity(key commons.Key) *entity {
	return &entity{
		key:  key,
		name: key.Code(),
	}
}

type entity struct {
	key         commons.Key
	name        string
	description string
}

func (e *entity) Key() commons.Key {
	if e == nil {
		return nil
	}

	return e.key
}

func (e *entity) Name() string {
	if e == nil {
		return ""
	}

	return e.name
}

func (e *entity) Description() string {
	if e == nil {
		return ""
	}

	return e.description
}
