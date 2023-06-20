package inventory

import (
	"github.com/shono-io/shono/commons"
)

type Event interface {
	Node
	Concept() commons.Reference
}

type EventSpec struct {
	NodeSpec
	Concept commons.Reference
}

type event struct {
	Spec EventSpec
}

func (e *event) Code() string {
	return e.Spec.Code
}

func (e *event) Summary() string {
	return e.Spec.Summary
}

func (e *event) Docs() string {
	return e.Spec.Docs
}

func (e *event) Status() commons.Status {
	return e.Spec.Status
}

func (e *event) Concept() commons.Reference {
	return e.Spec.Concept
}

func (e *event) Reference() commons.Reference {
	return e.Concept().Child("events", e.Code())
}
