package inventory

import "github.com/shono-io/shono/commons"

func NewEvent(scopeCode, conceptCode, code string) *EventBuilder {
	return &EventBuilder{
		event: &Event{
			Node: Node{
				Code: code,
			},
			Concept: NewConceptReference(scopeCode, conceptCode),
		},
	}
}

type EventBuilder struct {
	event *Event
}

func (b *EventBuilder) Summary(summary string) *EventBuilder {
	b.event.Summary = summary
	return b
}

func (b *EventBuilder) Docs(docs string) *EventBuilder {
	b.event.Docs = docs
	return b
}

func (b *EventBuilder) Status(status commons.Status) *EventBuilder {
	b.event.Status = status
	return b
}

func (b *EventBuilder) Build() *Event {
	return b.event
}
