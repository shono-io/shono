package internal

import (
	"github.com/shono-io/shono/core"
)

type InjectorSpec struct {
	NodeSpec
	Scope              core.Reference
	SourceSystem       core.Reference
	SourceSystemConfig map[string]any
	OutputEvents       []core.Reference
	Logic              core.Logic
}

type Injector struct {
	Spec InjectorSpec
}

func (i *Injector) Code() string {
	return i.Spec.Code
}

func (i *Injector) Summary() string {
	return i.Spec.Summary
}

func (i *Injector) Docs() string {
	return i.Spec.Docs
}

func (i *Injector) Status() core.Status {
	return i.Spec.Status
}

func (i *Injector) Scope() core.Reference {
	return i.Spec.Scope
}

func (i *Injector) SourceSystem() core.Reference {
	return i.Spec.SourceSystem
}

func (i *Injector) OutputEvents() []core.Reference {
	return i.Spec.OutputEvents
}

func (i *Injector) Logic() core.Logic {
	return i.Spec.Logic
}

func (i *Injector) SourceSystemConfig() map[string]any {
	return i.Spec.SourceSystemConfig
}

func (i *Injector) Reference() core.Reference {
	return i.Scope().Child("injectors", i.Code())
}
