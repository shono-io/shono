package shono

type Entity interface {
	FQN() string
	Code() string
	Name() string
	Description() string
}

func NewEntity(fqn, code, name, description string) Entity {
	return &entity{
		fqn:         fqn,
		code:        code,
		name:        name,
		description: description,
	}
}

func newEntity(fqn, code string) *entity {
	return &entity{
		fqn:  fqn,
		code: code,
		name: code,
	}
}

type entity struct {
	fqn         string
	code        string
	name        string
	description string
}

func (e *entity) FQN() string {
	if e == nil {
		return ""
	}

	return e.fqn
}

func (e *entity) Code() string {
	if e == nil {
		return ""
	}

	return e.code
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
