package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewInjector(scopeCode, code string) *InjectorBuilder {
	return &InjectorBuilder{
		injector: &Injector{
			Node: Node{
				Code: code,
			},
			Scope: NewScopeReference(scopeCode),
		},
	}
}

type InjectorBuilder struct {
	injector *Injector
}

func (b *InjectorBuilder) Summary(summary string) *InjectorBuilder {
	b.injector.Summary = summary
	return b
}

func (b *InjectorBuilder) Docs(docs string) *InjectorBuilder {
	b.injector.Docs = docs
	return b
}

func (b *InjectorBuilder) Status(status commons.Status) *InjectorBuilder {
	b.injector.Status = status
	return b
}

func (b *InjectorBuilder) Scope(code string) *InjectorBuilder {
	b.injector.Scope = commons.NewReference("scopes", code)
	return b
}

func (b *InjectorBuilder) Input(input Input) *InjectorBuilder {
	b.injector.Input = input
	return b
}

func (b *InjectorBuilder) OutputEvent(scopeCode, conceptCode, eventCode string) *InjectorBuilder {
	b.injector.OutputEvents = append(b.injector.OutputEvents, commons.NewReference("scopes", scopeCode).Child("concepts", conceptCode).Child("events", eventCode))
	return b
}

func (b *InjectorBuilder) Logic(lb *LogicBuilder) *InjectorBuilder {
	b.injector.Logic = lb.Build()
	return b
}

func (b *InjectorBuilder) Build() *Injector {
	return b.injector
}
