package inventory

import (
	"github.com/shono-io/shono/commons"
)

type ScopeSpec struct {
	NodeSpec
}

type Scope interface {
	Node
}

type scope struct {
	Spec ScopeSpec
}

func (s *scope) Code() string {
	return s.Spec.Code
}

func (s *scope) Summary() string {
	return s.Spec.Summary
}

func (s *scope) Docs() string {
	return s.Spec.Docs
}

func (s *scope) Status() commons.Status {
	return s.Spec.Status
}

func (s *scope) Reference() commons.Reference {
	return commons.NewReference("scopes", s.Code())
}
