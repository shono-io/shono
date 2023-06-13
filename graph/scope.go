package graph

import "github.com/benthosdev/benthos/v4/public/service"

var ScopeSpec = service.NewConfigSpec().
	Beta().
	Categories("Core").
	Summary("A scope is a unit within an organization.").
	Description("Scopes are used to group related concepts together. You aren't wrong thinking scopes sound a lot like bounded contexts. For example, the HR scope could group the concepts of employees, departments, and salaries.").
	Field(service.NewStringField("code").Description("The code of the scope. This is used to uniquely identify the scope.")).
	Field(service.NewStringField("name").Description("A human readable name for the scope")).
	Field(service.NewStringField("description").Description("Additional information about the scope"))

func NewScope(code, name, description string) Scope {
	return Scope{
		Code:        code,
		Name:        name,
		Description: description,
	}
}

type Scope struct {
	Code        string `yaml:"code"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type ScopeRepo interface {
	GetScope(code string) (*Scope, error)
	ListScopes() ([]Scope, error)
}
