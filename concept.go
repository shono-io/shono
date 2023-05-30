package shono

import (
	"fmt"
)

// == ENTITY ==========================================================================================================

type Concept interface {
	Entity
	ScopeCode() string
}

func NewConcept(scopeCode, code, name, description string) Concept {
	return &concept{
		scopeCode,
		NewEntity(fmt.Sprintf("%s:%s", scopeCode, code), code, name, description),
	}
}

type concept struct {
	scopeCode string
	Entity
}

func (c *concept) ScopeCode() string {
	return c.scopeCode
}
