package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewConceptReference(scopeCode, code string) commons.Reference {
	return NewScopeReference(scopeCode).Child("concepts", code)
}

type Concept struct {
	Node   `yaml:",inline"`
	Scope  commons.Reference
	Stored bool
}

func (c *Concept) Reference() commons.Reference {
	return NewConceptReference(c.Scope.Code(), c.Code)
}
