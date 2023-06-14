package core

import "github.com/shono-io/shono/internal"

func NewConcept(code string) *ConceptBuilder {
	return &ConceptBuilder{
		spec: internal.ConceptSpec{
			NodeSpec: internal.NodeSpec{
				Code: code,
			},
		},
	}
}

type ConceptBuilder struct {
	spec internal.ConceptSpec
}

func (b *ConceptBuilder) Summary(summary string) *ConceptBuilder {
	b.spec.Summary = summary
	return b
}

func (b *ConceptBuilder) Docs(docs string) *ConceptBuilder {
	b.spec.Docs = docs
	return b
}

func (b *ConceptBuilder) Status(status Status) *ConceptBuilder {
	b.spec.Status = status
	return b
}

func (b *ConceptBuilder) Scope(code string) *ConceptBuilder {
	b.spec.Scope = NewReference("scopes", code)
	return b
}

func (b *ConceptBuilder) Build() Concept {
	return &internal.Concept{Spec: b.spec}
}

type Concept interface {
	Node
	Scope() Reference
}
