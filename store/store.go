package store

import "fmt"

type store struct {
	scopeCode   string
	conceptCode string
	code        string
	name        string
	description string
}

func (e *store) ScopeCode() string {
	if e == nil {
		return ""
	}

	return e.scopeCode
}

func (e *store) ConceptCode() string {
	if e == nil {
		return ""
	}

	return e.conceptCode
}

func (e *store) FQN() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s:%s", e.scopeCode, e.conceptCode, e.code)
}

func (e *store) Code() string {
	if e == nil {
		return ""
	}

	return e.code
}

func (e *store) Name() string {
	if e == nil {
		return ""
	}

	return e.name
}

func (e *store) Description() string {
	if e == nil {
		return ""
	}

	return e.description
}
