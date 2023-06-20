package inventory

import (
	"github.com/shono-io/shono/commons"
)

type Reactor interface {
	Executable
	Concept() commons.Reference

	InputEvent() commons.Reference
	OutputEventCodes() []string
}

type ReactorSpec struct {
	NodeSpec
	Concept          commons.Reference
	InputEvent       commons.Reference
	OutputEventCodes []string
	Logic            Logic
}

type reactor struct {
	Spec ReactorSpec
}

func (r *reactor) Code() string {
	return r.Spec.Code
}

func (r *reactor) Summary() string {
	return r.Spec.Summary
}

func (r *reactor) Docs() string {
	return r.Spec.Docs
}

func (r *reactor) Status() commons.Status {
	return r.Spec.Status
}

func (r *reactor) Concept() commons.Reference {
	return r.Spec.Concept
}

func (r *reactor) InputEvent() commons.Reference {
	return r.Spec.InputEvent
}

func (r *reactor) OutputEventCodes() []string {
	return r.Spec.OutputEventCodes
}

func (r *reactor) Logic() Logic {
	return r.Spec.Logic
}

func (r *reactor) Reference() commons.Reference {
	return r.Concept().Child("reactors", r.Code())
}
