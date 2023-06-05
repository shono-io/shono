package graph

import "fmt"

type ConceptOpt func(*Concept)

func WithConceptName(name string) ConceptOpt {
	return func(c *Concept) {
		c.name = name
	}
}

func WithConceptDescription(description string) ConceptOpt {
	return func(c *Concept) {
		c.description = description
	}
}

func WithConceptPluralName(plural string) ConceptOpt {
	return func(c *Concept) {
		c.plural = plural
	}
}

func WithConceptSingleName(single string) ConceptOpt {
	return func(c *Concept) {
		c.single = single
	}
}

func NewConcept(key Key, opts ...ConceptOpt) Concept {
	result := Concept{
		key.Parent(),
		key,
		key.Code(),
		"",
		fmt.Sprintf("%ss", key.Code()),
		key.Code(),
	}

	for _, opt := range opts {
		opt(&result)
	}

	return result
}

type Concept struct {
	scopeKey    Key
	key         Key
	name        string
	description string
	plural      string
	single      string
}

func (c Concept) Key() Key {
	return c.key
}

func (c Concept) Name() string {
	return c.name
}

func (c Concept) Description() string {
	return c.description
}

func (c Concept) ScopeKey() Key {
	return c.scopeKey
}

func (c Concept) Plural() string {
	return c.plural
}

func (c Concept) Single() string {
	return c.single
}
