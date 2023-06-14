package core

import "github.com/shono-io/shono/internal"

func NewScope(code string) *ScopeBuilder {
	return &ScopeBuilder{
		spec: internal.ScopeSpec{
			NodeSpec: internal.NodeSpec{
				Code: code,
			},
		},
	}
}

type ScopeBuilder struct {
	spec internal.ScopeSpec
}

func (s *ScopeBuilder) Summary(summary string) *ScopeBuilder {
	s.spec.Summary = summary
	return s
}

func (s *ScopeBuilder) Docs(docs string) *ScopeBuilder {
	s.spec.Docs = docs
	return s
}

func (s *ScopeBuilder) Status(status Status) *ScopeBuilder {
	s.spec.Status = status
	return s
}

func (s *ScopeBuilder) Build() Scope {
	return &internal.Scope{Spec: s.spec}
}

type Scope interface {
	Node
}
