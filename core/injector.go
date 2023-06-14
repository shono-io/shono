package core

import "github.com/shono-io/shono/internal"

func NewInjector(code string) *InjectorBuilder {
	return &InjectorBuilder{
		spec: internal.InjectorSpec{
			NodeSpec: internal.NodeSpec{
				Code: code,
			},
		},
	}
}

type InjectorBuilder struct {
	spec internal.InjectorSpec
}

func (b *InjectorBuilder) Summary(summary string) *InjectorBuilder {
	b.spec.Summary = summary
	return b
}

func (b *InjectorBuilder) Docs(docs string) *InjectorBuilder {
	b.spec.Docs = docs
	return b
}

func (b *InjectorBuilder) Status(status Status) *InjectorBuilder {
	b.spec.Status = status
	return b
}

func (b *InjectorBuilder) Scope(code string) *InjectorBuilder {
	b.spec.Scope = NewReference("scopes", code)
	return b
}

func (b *InjectorBuilder) SourceSystem(code string, config map[string]any) *InjectorBuilder {
	b.spec.SourceSystem = NewReference("systems", code)
	b.spec.SourceSystemConfig = config
	return b
}

func (b *InjectorBuilder) OutputEvent(scopeCode, conceptCode, eventCode string) *InjectorBuilder {
	b.spec.OutputEvents = append(b.spec.OutputEvents, NewReference("scopes", scopeCode).Child("concepts", conceptCode).Child("events", eventCode))
	return b
}

func (b *InjectorBuilder) Logic(lb LogicBuilder) *InjectorBuilder {
	b.spec.Logic = lb.Build()
	return b
}

func (b *InjectorBuilder) Build() Injector {
	return &internal.Injector{Spec: b.spec}
}

type Injector interface {
	Executable
	Scope() Reference

	SourceSystem() Reference
	SourceSystemConfig() map[string]any
	OutputEvents() []Reference
}
