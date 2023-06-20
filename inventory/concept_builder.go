package inventory

import "github.com/shono-io/shono/commons"

func NewConceptReference(scopeCode, code string) commons.Reference {
	return NewScopeReference(scopeCode).Child("concepts", code)
}

func NewConcept(scopeCode, code string) *ConceptBuilder {
	return &ConceptBuilder{
		spec: ConceptSpec{
			NodeSpec: NodeSpec{
				Code: code,
			},
			Scope: NewScopeReference(scopeCode),
		},
	}
}

type ConceptBuilder struct {
	spec ConceptSpec
}

func (b *ConceptBuilder) Summary(summary string) *ConceptBuilder {
	b.spec.Summary = summary
	return b
}

func (b *ConceptBuilder) Docs(docs string) *ConceptBuilder {
	b.spec.Docs = docs
	return b
}

func (b *ConceptBuilder) Status(status commons.Status) *ConceptBuilder {
	b.spec.Status = status
	return b
}

func (b *ConceptBuilder) Build() Concept {
	return &concept{Spec: b.spec}
}
