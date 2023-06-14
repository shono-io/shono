package core

import "github.com/shono-io/shono/internal"

func NewReactor(code string) *ReaktorBuilder {
	return &ReaktorBuilder{
		spec: internal.ReactorSpec{
			NodeSpec: internal.NodeSpec{
				Code: code,
			},
		},
	}
}

type ReaktorBuilder struct {
	spec internal.ReactorSpec
}

func (r *ReaktorBuilder) Summary(summary string) *ReaktorBuilder {
	r.spec.Summary = summary
	return r
}

func (r *ReaktorBuilder) Docs(docs string) *ReaktorBuilder {
	r.spec.Docs = docs
	return r
}

func (r *ReaktorBuilder) Status(status Status) *ReaktorBuilder {
	r.spec.Status = status
	return r
}

func (r *ReaktorBuilder) Concept(scopeCode, conceptCode string) *ReaktorBuilder {
	r.spec.Concept = NewReference("scopes", scopeCode).Child("concepts", conceptCode)
	return r
}

func (r *ReaktorBuilder) InputEvent(scopeCode, conceptCode, eventCode string) *ReaktorBuilder {
	r.spec.InputEvent = NewReference("scopes", scopeCode).Child("concepts", conceptCode).Child("events", eventCode)
	return r
}

func (r *ReaktorBuilder) OutputEvent(scopeCode, conceptCode, eventCode string) *ReaktorBuilder {
	r.spec.OutputEvents = append(r.spec.OutputEvents, NewReference("scopes", scopeCode).Child("concepts", conceptCode).Child("events", eventCode))
	return r
}

func (r *ReaktorBuilder) Logic(b LogicBuilder) *ReaktorBuilder {
	r.spec.Logic = b.Build()
	return r
}

func (r *ReaktorBuilder) Build() Reactor {
	return &internal.Reactor{Spec: r.spec}
}

type Reactor interface {
	Executable
	Concept() Reference

	InputEvent() Reference
	OutputEvents() []Reference
}
