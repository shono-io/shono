package internal

import (
	"github.com/shono-io/shono/core"
)

type ReactorSpec struct {
	NodeSpec
	Concept      core.Reference
	InputEvent   core.Reference
	OutputEvents []core.Reference
	Logic        core.Logic
}

type Reactor struct {
	Spec ReactorSpec
}

func (r *Reactor) Code() string {
	return r.Spec.Code
}

func (r *Reactor) Summary() string {
	return r.Spec.Summary
}

func (r *Reactor) Docs() string {
	return r.Spec.Docs
}

func (r *Reactor) Status() core.Status {
	return r.Spec.Status
}

func (r *Reactor) Concept() core.Reference {
	return r.Spec.Concept
}

func (r *Reactor) InputEvent() core.Reference {
	return r.Spec.InputEvent
}

func (r *Reactor) OutputEvents() []core.Reference {
	return r.Spec.OutputEvents
}

func (r *Reactor) Logic() core.Logic {
	return r.Spec.Logic
}

func (r *Reactor) Reference() core.Reference {
	return r.Concept().Child("reactors", r.Code())
}
