package inventory

import (
	"github.com/shono-io/shono/commons"
)

func NewExtractorReference(scopeCode, extractorCode string) commons.Reference {
	return NewScopeReference(scopeCode).Child("extractors", extractorCode)
}

type Extractor struct {
	Node
	Scope       commons.Reference
	Output      Output
	InputEvents []commons.Reference
	Logic       Logic
}

func (e *Extractor) Reference() commons.Reference {
	return NewExtractorReference(e.Scope.Code(), e.Code)
}
