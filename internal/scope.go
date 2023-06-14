package internal

import (
	"github.com/shono-io/shono/core"
)

type ScopeSpec struct {
	NodeSpec
}

type Scope struct {
	Spec ScopeSpec
}

func (s *Scope) Code() string {
	return s.Spec.Code
}

func (s *Scope) Summary() string {
	return s.Spec.Summary
}

func (s *Scope) Docs() string {
	return s.Spec.Docs
}

func (s *Scope) Status() core.Status {
	return s.Spec.Status
}

func (s *Scope) Reference() core.Reference {
	return core.NewReference("scopes", s.Code())
}
