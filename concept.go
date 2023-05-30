package shono

import (
	"fmt"
)

type Concept interface {
	Entity
	ScopeCode() string
	Event(code string, opts ...EventOpt) Event
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

func NewConcept(scopeCode, code string, opts ...ConceptOpt) Concept {
	result := &concept{
		scopeCode,
		newEntity(fmt.Sprintf("%s:%s", scopeCode, code), code),
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type concept struct {
	scopeCode string
	*entity
}

func (c *concept) ScopeCode() string {
	return c.scopeCode
}

func (c *concept) Event(code string, opts ...EventOpt) Event {
	return NewEvent(c.scopeCode, c.Code(), code, opts...)
}
