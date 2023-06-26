package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewScope(code string) *ScopeBuilder {
	return &ScopeBuilder{
		scope: &Scope{
			Node: Node{
				Code: code,
			},
		},
	}
}

type ScopeBuilder struct {
	scope *Scope
}

func (s *ScopeBuilder) Summary(summary string) *ScopeBuilder {
	s.scope.Summary = summary
	return s
}

func (s *ScopeBuilder) Docs(docs string) *ScopeBuilder {
	s.scope.Docs = docs
	return s
}

func (s *ScopeBuilder) Status(status commons.Status) *ScopeBuilder {
	s.scope.Status = status
	return s
}

func (s *ScopeBuilder) Build() *Scope {
	return s.scope
}
