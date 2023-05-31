package shono

type Entity interface {
	Key() Key
	Name() string
	Description() string
}

func NewEntity(key Key, name, description string) Entity {
	return &entity{
		key:         key,
		name:        name,
		description: description,
	}
}

func newEntity(key Key) *entity {
	return &entity{
		key:  key,
		name: key.Code(),
	}
}

type entity struct {
	key         Key
	name        string
	description string
}

func (e *entity) Key() Key {
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
