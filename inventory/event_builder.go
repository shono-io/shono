package inventory

import "github.com/shono-io/shono/commons"

func NewEventReference(scopeCode, conceptCode, code string) commons.Reference {
	return NewConceptReference(scopeCode, conceptCode).Child("events", code)
}

func NewEvent(scopeCode, conceptCode, code string) *EventBuilder {
	return &EventBuilder{
		spec: EventSpec{
			NodeSpec: NodeSpec{
				Code: code,
			},
			Concept: NewConceptReference(scopeCode, conceptCode),
		},
	}
}

type EventBuilder struct {
	spec EventSpec
}

func (b *EventBuilder) Summary(summary string) *EventBuilder {
	b.spec.Summary = summary
	return b
}

func (b *EventBuilder) Docs(docs string) *EventBuilder {
	b.spec.Docs = docs
	return b
}

func (b *EventBuilder) Status(status commons.Status) *EventBuilder {
	b.spec.Status = status
	return b
}

func (b *EventBuilder) Build() Event {
	return &event{Spec: b.spec}
}
