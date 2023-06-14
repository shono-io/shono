package graph

import "github.com/shono-io/shono/core"

//var ScopeSpec = service.NewConfigSpec().
//	Beta().
//	Categories("Core").
//	Summary("A scope is a unit within an organization.").
//	Description("Scopes are used to group related concepts together. You aren't wrong thinking scopes sound a lot like bounded contexts. For example, the HR scope could group the concepts of employees, departments, and salaries.").
//	Field(service.NewStringField("code").Description("The code of the scope. This is used to uniquely identify the scope.")).
//	Field(service.NewStringField("name").Description("A human readable name for the scope")).
//	Field(service.NewStringField("description").Description("Additional information about the scope"))

type ScopeRepo interface {
	GetScope(code string) (*Scope, error)
	ListScopes() ([]Scope, error)
}

func NewScopeSpec(code string) *ScopeSpec {
	return &ScopeSpec{scope: &Scope{Code: code}}
}

type ScopeSpec struct {
	scope *Scope
}

type Scope struct {
	Code    string `yaml:"code"`
	Summary string `yaml:"summary"`
	Docs    string `yaml:"docs"`
	Status  Status `yaml:"status"`

	Concepts []core.Concept `yaml:"concepts"`
}

func (s *ScopeSpec) Summary(summary string) *ScopeSpec {
	s.scope.Summary = summary
	return s
}

func (s *ScopeSpec) Docs(docs string) *ScopeSpec {
	s.scope.Docs = docs
	return s
}

func (s *ScopeSpec) Stable() *ScopeSpec {
	s.scope.Status = StatusStable
	return s
}

func (s *ScopeSpec) Beta() *ScopeSpec {
	s.scope.Status = StatusBeta
	return s
}

func (s *ScopeSpec) Experimental() *ScopeSpec {
	s.scope.Status = StatusExperimental
	return s
}

func (s *ScopeSpec) Deprecated() *ScopeSpec {
	s.scope.Status = StatusDeprecated
	return s
}

func (s *ScopeSpec) Concept(concepts ...*ConceptSpec) *ScopeSpec {
	for _, concept := range concepts {
		s.scope.Concepts = append(s.scope.Concepts, *concept.concept)
	}

	return s
}
