package inventory

import "github.com/shono-io/shono/commons"

func NewScopeReference(code string) commons.Reference {
	return commons.NewReference("scopes", code)
}

type Scope struct {
	Node
}

func (s *Scope) Reference() commons.Reference {
	return NewScopeReference(s.Code)
}
