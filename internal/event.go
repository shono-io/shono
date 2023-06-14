package internal

import (
	"github.com/shono-io/shono/core"
)

type EventSpec struct {
	NodeSpec
	Concept core.Reference
}

type Event struct {
	Spec EventSpec
}

func (e *Event) Code() string {
	return e.Spec.Code
}

func (e *Event) Summary() string {
	return e.Spec.Summary
}

func (e *Event) Docs() string {
	return e.Spec.Docs
}

func (e *Event) Status() core.Status {
	return e.Spec.Status
}

func (e *Event) Concept() core.Reference {
	return e.Spec.Concept
}

func (e *Event) Reference() core.Reference {
	return e.Concept().Child("events", e.Code())
}
