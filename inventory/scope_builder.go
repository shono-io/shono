package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewScopeReference(code string) commons.Reference {
	return commons.NewReference("scopes", code)
}

func NewScope(code string) *ScopeBuilder {
	return &ScopeBuilder{
		spec: ScopeSpec{
			NodeSpec: NodeSpec{
				Code: code,
			},
		},
	}
}

type ScopeBuilder struct {
	spec ScopeSpec
}

func (s *ScopeBuilder) Summary(summary string) *ScopeBuilder {
	s.spec.Summary = summary
	return s
}

func (s *ScopeBuilder) Docs(docs string) *ScopeBuilder {
	s.spec.Docs = docs
	return s
}

func (s *ScopeBuilder) Status(status commons.Status) *ScopeBuilder {
	s.spec.Status = status
	return s
}

func (s *ScopeBuilder) Build() Scope {
	return &scope{Spec: s.spec}
}
