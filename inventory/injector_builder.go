package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewInjectorReference(scopeCode, injectorCode string) commons.Reference {
	return NewScopeReference(scopeCode).Child("injectors", injectorCode)
}

func NewInjector(scopeCode, code string) *InjectorBuilder {
	return &InjectorBuilder{
		spec: InjectorSpec{
			NodeSpec: NodeSpec{
				Code: code,
			},
			Scope: NewScopeReference(scopeCode),
		},
	}
}

type InjectorBuilder struct {
	spec InjectorSpec
}

func (b *InjectorBuilder) Summary(summary string) *InjectorBuilder {
	b.spec.Summary = summary
	return b
}

func (b *InjectorBuilder) Docs(docs string) *InjectorBuilder {
	b.spec.Docs = docs
	return b
}

func (b *InjectorBuilder) Status(status commons.Status) *InjectorBuilder {
	b.spec.Status = status
	return b
}

func (b *InjectorBuilder) Scope(code string) *InjectorBuilder {
	b.spec.Scope = commons.NewReference("scopes", code)
	return b
}

func (b *InjectorBuilder) Input(input Input) *InjectorBuilder {
	b.spec.Input = input
	return b
}

func (b *InjectorBuilder) OutputEvent(scopeCode, conceptCode, eventCode string) *InjectorBuilder {
	b.spec.OutputEvents = append(b.spec.OutputEvents, commons.NewReference("scopes", scopeCode).Child("concepts", conceptCode).Child("events", eventCode))
	return b
}

func (b *InjectorBuilder) Logic(lb LogicBuilder) *InjectorBuilder {
	b.spec.Logic = lb.Build()
	return b
}

func (b *InjectorBuilder) Build() Injector {
	return &injector{Spec: b.spec}
}
