package inventory

import (
	"github.com/shono-io/shono/commons"
)

type Injector interface {
	Executable
	Scope() commons.Reference

	Input() Input
	OutputEvents() []commons.Reference
}

type InjectorSpec struct {
	NodeSpec
	Scope        commons.Reference
	Input        Input
	OutputEvents []commons.Reference
	Logic        Logic
}

type injector struct {
	Spec InjectorSpec
}

func (i *injector) Code() string {
	return i.Spec.Code
}

func (i *injector) Summary() string {
	return i.Spec.Summary
}

func (i *injector) Docs() string {
	return i.Spec.Docs
}

func (i *injector) Status() commons.Status {
	return i.Spec.Status
}

func (i *injector) Scope() commons.Reference {
	return i.Spec.Scope
}

func (i *injector) Input() Input {
	return i.Spec.Input
}

func (i *injector) OutputEvents() []commons.Reference {
	return i.Spec.OutputEvents
}

func (i *injector) Logic() Logic {
	return i.Spec.Logic
}

func (i *injector) Reference() commons.Reference {
	return i.Scope().Child("injectors", i.Code())
}
