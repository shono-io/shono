package shono

import (
	"fmt"
	"strings"
)

type Concept interface {
	Entity
	Event(code string, opts ...EventOpt) Event
	Plural() string
	Single() string
}

type ConceptOpt func(*concept)

func WithConceptName(name string) ConceptOpt {
	return func(c *concept) {
		c.entity.name = name
	}
}

func WithConceptDescription(description string) ConceptOpt {
	return func(c *concept) {
		c.entity.description = description
	}
}

func WithConceptPluralName(plural string) ConceptOpt {
	return func(c *concept) {
		c.plural = plural
	}
}

func WithConceptSingleName(single string) ConceptOpt {
	return func(c *concept) {
		c.single = single
	}
}

func NewConcept(scopeKey Key, code string, opts ...ConceptOpt) Concept {
	result := &concept{
		newEntity(scopeKey.Child("concept", code)),
		fmt.Sprintf("%ss", strings.ToTitle(code)),
		strings.ToTitle(code),
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type concept struct {
	*entity
	plural string
	single string
}

func (c *concept) Event(code string, opts ...EventOpt) Event {
	return NewEvent(c.Key(), code, opts...)
}

func (c *concept) Plural() string {
	return c.plural
}

func (c *concept) Single() string {
	return c.single
}
