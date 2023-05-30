package shono

type Scope interface {
	Entity
}

func NewScope(code, name, description string) Scope {
	return &scope{
		NewEntity(code, code, name, description),
	}
}

type scope struct {
	Entity
}
