package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewExtractor(scopeCode, code string) *ExtractorBuilder {
	return &ExtractorBuilder{
		spec: ExtractorSpec{
			NodeSpec: NodeSpec{
				Code: code,
			},
			Scope: NewScopeReference(scopeCode),
		},
	}
}

type ExtractorBuilder struct {
	spec ExtractorSpec
}

func (e *ExtractorBuilder) Summary(summary string) *ExtractorBuilder {
	e.spec.Summary = summary
	return e
}

func (e *ExtractorBuilder) Docs(docs string) *ExtractorBuilder {
	e.spec.Docs = docs
	return e
}

func (e *ExtractorBuilder) Status(status commons.Status) *ExtractorBuilder {
	e.spec.Status = status
	return e
}

func (e *ExtractorBuilder) Scope(code string) *ExtractorBuilder {
	e.spec.Scope = commons.NewReference("scopes", code)
	return e
}

func (e *ExtractorBuilder) Output(output Output) *ExtractorBuilder {
	e.spec.Output = output
	return e
}

func (e *ExtractorBuilder) InputEvent(scopeCode, conceptCode, eventCode string) *ExtractorBuilder {
	e.spec.InputEvents = append(e.spec.InputEvents, commons.NewReference("scopes", scopeCode).Child("concepts", conceptCode).Child("events", eventCode))
	return e
}

func (e *ExtractorBuilder) Logic(b LogicBuilder) *ExtractorBuilder {
	e.spec.Logic = b.Build()
	return e
}

func (e *ExtractorBuilder) Build() Extractor {
	return &extractor{Spec: e.spec}
}
