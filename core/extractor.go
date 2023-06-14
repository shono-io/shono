package core

import "github.com/shono-io/shono/internal"

func NewExtractor(code string) *ExtractorBuilder {
	return &ExtractorBuilder{
		spec: internal.ExtractorSpec{
			NodeSpec: internal.NodeSpec{
				Code: code,
			},
		},
	}
}

type ExtractorBuilder struct {
	spec internal.ExtractorSpec
}

func (e *ExtractorBuilder) Summary(summary string) *ExtractorBuilder {
	e.spec.Summary = summary
	return e
}

func (e *ExtractorBuilder) Docs(docs string) *ExtractorBuilder {
	e.spec.Docs = docs
	return e
}

func (e *ExtractorBuilder) Status(status Status) *ExtractorBuilder {
	e.spec.Status = status
	return e
}

func (e *ExtractorBuilder) Scope(code string) *ExtractorBuilder {
	e.spec.Scope = NewReference("scopes", code)
	return e
}

func (e *ExtractorBuilder) TargetSystem(code string, config map[string]any) *ExtractorBuilder {
	e.spec.TargetSystem = NewReference("systems", code)
	e.spec.TargetSystemConfig = config
	return e
}

func (e *ExtractorBuilder) InputEvent(scopeCode, conceptCode, eventCode string) *ExtractorBuilder {
	e.spec.InputEvents = append(e.spec.InputEvents, NewReference("scopes", scopeCode).Child("concepts", conceptCode).Child("events", eventCode))
	return e
}

func (e *ExtractorBuilder) Logic(b LogicBuilder) *ExtractorBuilder {
	e.spec.Logic = b.Build()
	return e
}

func (e *ExtractorBuilder) Build() Extractor {
	return &internal.Extractor{Spec: e.spec}
}

type Extractor interface {
	Executable
	Scope() Reference

	InputEvents() []Reference
	TargetSystem() Reference
	TargetSystemConfig() map[string]any
}
