package core

import "github.com/shono-io/shono/internal"

func NewEvent(code string) *EventBuilder {
	return &EventBuilder{
		spec: internal.EventSpec{
			NodeSpec: internal.NodeSpec{
				Code: code,
			},
		},
	}
}

type EventBuilder struct {
	spec internal.EventSpec
}

func (b *EventBuilder) Summary(summary string) *EventBuilder {
	b.spec.Summary = summary
	return b
}

func (b *EventBuilder) Docs(docs string) *EventBuilder {
	b.spec.Docs = docs
	return b
}

func (b *EventBuilder) Status(status Status) *EventBuilder {
	b.spec.Status = status
	return b
}

func (b *EventBuilder) Concept(scopeCode, conceptCode string) *EventBuilder {
	b.spec.Concept = NewReference("scopes", scopeCode).Child("concepts", conceptCode)
	return b
}

func (b *EventBuilder) Build() Event {
	return &internal.Event{Spec: b.spec}
}

type Event interface {
	Node
	Concept() Reference
}
