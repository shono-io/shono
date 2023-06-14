package internal

import (
	"github.com/shono-io/shono/core"
)

type ExtractorSpec struct {
	NodeSpec
	Scope              core.Reference
	TargetSystem       core.Reference
	TargetSystemConfig map[string]any
	InputEvents        []core.Reference
	Logic              core.Logic
}

type Extractor struct {
	Spec ExtractorSpec
}

func (e *Extractor) Code() string {
	return e.Spec.Code
}

func (e *Extractor) Summary() string {
	return e.Spec.Summary
}

func (e *Extractor) Docs() string {
	return e.Spec.Docs
}

func (e *Extractor) Status() core.Status {
	return e.Spec.Status
}

func (e *Extractor) Scope() core.Reference {
	return e.Spec.Scope
}

func (e *Extractor) TargetSystem() core.Reference {
	return e.Spec.TargetSystem
}

func (e *Extractor) InputEvents() []core.Reference {
	return e.Spec.InputEvents
}

func (e *Extractor) Logic() core.Logic {
	return e.Spec.Logic
}

func (e *Extractor) TargetSystemConfig() map[string]any {
	return e.Spec.TargetSystemConfig
}

func (e *Extractor) Reference() core.Reference {
	return e.Scope().Child("extractors", e.Code())
}
