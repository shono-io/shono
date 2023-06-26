package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewInjectorReference(scopeCode, injectorCode string) commons.Reference {
	return NewScopeReference(scopeCode).Child("injectors", injectorCode)
}

type Injector struct {
	Node
	Scope        commons.Reference
	Input        Input
	OutputEvents []commons.Reference
	Logic        Logic
}

func (i *Injector) Reference() commons.Reference {
	return NewInjectorReference(i.Scope.Code(), i.Code)
}
