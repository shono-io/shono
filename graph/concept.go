package graph

import (
	"fmt"
	"github.com/shono-io/shono/commons"
)

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

func WithRequest(request ...Request) ConceptOpt {
	return func(c *Concept) {
		c.requests = append(c.requests, request...)
	}
}

func NewConcept(key commons.Key, opts ...ConceptOpt) Concept {
	result := Concept{
		key.Parent(),
		key,
		key.Code(),
		"",
		fmt.Sprintf("%ss", key.Code()),
		key.Code(),
		[]Request{},
	}

	for _, opt := range opts {
		opt(&result)
	}

	return result
}

type Concept struct {
	scopeKey    commons.Key
	key         commons.Key
	name        string
	description string
	plural      string
	single      string
	requests    []Request
}

func (c Concept) Key() commons.Key {
	return c.key
}

func (c Concept) Name() string {
	return c.name
}

func (c Concept) Description() string {
	return c.description
}

func (c Concept) ScopeKey() commons.Key {
	return c.scopeKey
}

func (c Concept) Plural() string {
	return c.plural
}

func (c Concept) Single() string {
	return c.single
}

func (c Concept) Requests() []Request {
	return c.requests
}
