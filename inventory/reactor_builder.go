package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewReactor(scopeCode, conceptCode, code string) *ReaktorBuilder {
	return &ReaktorBuilder{
		reactor: &Reactor{
			Node: Node{
				Code: code,
			},
			Concept: NewConceptReference(scopeCode, conceptCode),
		},
	}
}

type ReaktorBuilder struct {
	reactor *Reactor
}

func (r *ReaktorBuilder) Summary(summary string) *ReaktorBuilder {
	r.reactor.Summary = summary
	return r
}

func (r *ReaktorBuilder) Docs(docs string) *ReaktorBuilder {
	r.reactor.Docs = docs
	return r
}

func (r *ReaktorBuilder) Status(status commons.Status) *ReaktorBuilder {
	r.reactor.Status = status
	return r
}

func (r *ReaktorBuilder) InputEvent(eventRef commons.Reference) *ReaktorBuilder {
	r.reactor.InputEvent = eventRef
	return r
}

func (r *ReaktorBuilder) OutputEventCodes(eventCodes ...string) *ReaktorBuilder {
	r.reactor.OutputEventCodes = eventCodes
	return r
}

func (r *ReaktorBuilder) Logic(b *LogicBuilder) *ReaktorBuilder {
	r.reactor.Logic = b.Build()
	return r
}

func (r *ReaktorBuilder) Build() *Reactor {
	return r.reactor
}
