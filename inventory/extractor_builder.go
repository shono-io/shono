package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewExtractor(scopeCode, code string) *ExtractorBuilder {
	return &ExtractorBuilder{
		extractor: &Extractor{
			Node: Node{
				Code: code,
			},
			Scope: NewScopeReference(scopeCode),
		},
	}
}

type ExtractorBuilder struct {
	extractor *Extractor
}

func (e *ExtractorBuilder) Summary(summary string) *ExtractorBuilder {
	e.extractor.Summary = summary
	return e
}

func (e *ExtractorBuilder) Docs(docs string) *ExtractorBuilder {
	e.extractor.Docs = docs
	return e
}

func (e *ExtractorBuilder) Status(status commons.Status) *ExtractorBuilder {
	e.extractor.Status = status
	return e
}

func (e *ExtractorBuilder) Scope(code string) *ExtractorBuilder {
	e.extractor.Scope = commons.NewReference("scopes", code)
	return e
}

func (e *ExtractorBuilder) Output(output Output) *ExtractorBuilder {
	e.extractor.Output = output
	return e
}

func (e *ExtractorBuilder) InputEvent(scopeCode, conceptCode, eventCode string) *ExtractorBuilder {
	e.extractor.InputEvents = append(e.extractor.InputEvents, commons.NewReference("scopes", scopeCode).Child("concepts", conceptCode).Child("events", eventCode))
	return e
}

func (e *ExtractorBuilder) Logic(b *LogicBuilder) *ExtractorBuilder {
	e.extractor.Logic = b.Build()
	return e
}

func (e *ExtractorBuilder) Build() *Extractor {
	return e.extractor
}
