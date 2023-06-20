package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewReactorReference(scopeCode, conceptCode, code string) commons.Reference {
	return NewConceptReference(scopeCode, conceptCode).Child("reactors", code)
}

func NewReactor(scopeCode, conceptCode, code string) *ReaktorBuilder {
	return &ReaktorBuilder{
		spec: ReactorSpec{
			NodeSpec: NodeSpec{
				Code: code,
			},
			Concept: NewConceptReference(scopeCode, conceptCode),
		},
	}
}

type ReaktorBuilder struct {
	spec ReactorSpec
}

func (r *ReaktorBuilder) Summary(summary string) *ReaktorBuilder {
	r.spec.Summary = summary
	return r
}

func (r *ReaktorBuilder) Docs(docs string) *ReaktorBuilder {
	r.spec.Docs = docs
	return r
}

func (r *ReaktorBuilder) Status(status commons.Status) *ReaktorBuilder {
	r.spec.Status = status
	return r
}

func (r *ReaktorBuilder) InputEvent(eventRef commons.Reference) *ReaktorBuilder {
	r.spec.InputEvent = eventRef
	return r
}

func (r *ReaktorBuilder) OutputEventCodes(eventCodes ...string) *ReaktorBuilder {
	r.spec.OutputEventCodes = eventCodes
	return r
}

func (r *ReaktorBuilder) Logic(b LogicBuilder) *ReaktorBuilder {
	r.spec.Logic = b.Build()
	return r
}

func (r *ReaktorBuilder) Build() Reactor {
	return &reactor{Spec: r.spec}
}
