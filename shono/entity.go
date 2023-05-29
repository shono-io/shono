package shono

type Entity interface {
	FQN() string
	GetCode() string
	GetName() string
	GetDescription() string
}

func NewEntity(fqn, code, name, description string) Entity {
	return &entity{
		fqn:         fqn,
		code:        code,
		name:        name,
		description: description,
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

func (e *entity) GetCode() string {
	if e == nil {
		return ""
	}

	return e.code
}

func (e *entity) GetName() string {
	if e == nil {
		return ""
	}

	return e.name
}

func (e *entity) GetDescription() string {
	if e == nil {
		return ""
	}

	return e.description
}
