package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewExtractorReference(scopeCode, extractorCode string) commons.Reference {
	return NewScopeReference(scopeCode).Child("extractors", extractorCode)
}

type Extractor interface {
	Executable
	Scope() commons.Reference

	InputEvents() []commons.Reference
	Output() Output
}

type ExtractorSpec struct {
	NodeSpec
	Scope       commons.Reference
	Output      Output
	InputEvents []commons.Reference
	Logic       Logic
}

type extractor struct {
	Spec ExtractorSpec
}

func (e *extractor) Code() string {
	return e.Spec.Code
}

func (e *extractor) Summary() string {
	return e.Spec.Summary
}

func (e *extractor) Docs() string {
	return e.Spec.Docs
}

func (e *extractor) Status() commons.Status {
	return e.Spec.Status
}

func (e *extractor) Scope() commons.Reference {
	return e.Spec.Scope
}

func (e *extractor) Output() Output {
	return e.Spec.Output
}

func (e *extractor) InputEvents() []commons.Reference {
	return e.Spec.InputEvents
}

func (e *extractor) Logic() Logic {
	return e.Spec.Logic
}

func (e *extractor) Reference() commons.Reference {
	return e.Scope().Child("extractors", e.Code())
}
