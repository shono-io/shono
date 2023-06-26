package inventory

import "github.com/shono-io/shono/commons"

func NewConcept(scopeCode, code string) *ConceptBuilder {
	return &ConceptBuilder{
		concept: &Concept{
			Node: Node{
				Code: code,
			},
			Scope: NewScopeReference(scopeCode),
		},
	}
}

type ConceptBuilder struct {
	concept *Concept
}

func (b *ConceptBuilder) Summary(summary string) *ConceptBuilder {
	b.concept.Summary = summary
	return b
}

func (b *ConceptBuilder) Docs(docs string) *ConceptBuilder {
	b.concept.Docs = docs
	return b
}

func (b *ConceptBuilder) Status(status commons.Status) *ConceptBuilder {
	b.concept.Status = status
	return b
}

func (b *ConceptBuilder) Stored() *ConceptBuilder {
	b.concept.Stored = true
	return b
}

func (b *ConceptBuilder) Build() *Concept {
	return b.concept
}
