package graph

import (
	"fmt"
	"strings"
)

func NewConceptSpec(code string) *ConceptSpec {
	return &ConceptSpec{
		concept: &Concept{
			Code:   code,
			Single: code,
			Plural: code + "s",
		},
	}
}

type ConceptSpec struct {
	concept *Concept
}

func (c *ConceptSpec) Summary(summary string) *ConceptSpec {
	c.concept.Summary = summary
	return c
}

func (c *ConceptSpec) Docs(docs string) *ConceptSpec {
	c.concept.Docs = docs
	return c
}

func (c *ConceptSpec) Stable() *ConceptSpec {
	c.concept.Status = StatusStable
	return c
}

func (c *ConceptSpec) Beta() *ConceptSpec {
	c.concept.Status = StatusBeta
	return c
}

func (c *ConceptSpec) Experimental() *ConceptSpec {
	c.concept.Status = StatusExperimental
	return c
}

func (c *ConceptSpec) Deprecated() *ConceptSpec {
	c.concept.Status = StatusDeprecated
	return c
}

func (c *ConceptSpec) Store(storageKey, collection string) *ConceptSpec {
	c.concept.Store = &ConceptStore{
		StorageKey: storageKey,
		Collection: collection,
	}

	return c
}

func (c *ConceptSpec) Single(single string) *ConceptSpec {
	c.concept.Single = single
	return c
}

func (c *ConceptSpec) Plural(plural string) *ConceptSpec {
	c.concept.Plural = plural
	return c
}

func (c *ConceptSpec) Event(events ...*EventSpec) *ConceptSpec {
	for _, event := range events {
		c.concept.Events = append(c.concept.Events, *event.event)
	}

	return c
}

type Concept struct {
	Code    string        `yaml:"code"`
	Status  Status        `yaml:"status"`
	Summary string        `yaml:"summary,omitempty"`
	Docs    string        `yaml:"docs,omitempty"`
	Plural  string        `yaml:"plural,omitempty"`
	Single  string        `yaml:"single,omitempty"`
	Store   *ConceptStore `yaml:"store,omitempty"`

	Events   []Event   `yaml:"events,omitempty"`
	Reaktors []Reaktor `yaml:"reaktors,omitempty"`
}

type ConceptStore struct {
	StorageKey string `yaml:"storage"`
	Collection string `yaml:"collection"`
}

func ParseConceptReference(input string) (ConceptReference, error) {
	parts := strings.Split(input, "__")
	if len(parts) != 2 {
		return ConceptReference{}, fmt.Errorf("invalid concept reference: %q", input)
	}

	return ConceptReference{
		ScopeCode: parts[0],
		Code:      parts[1],
	}, nil
}

type ConceptReference struct {
	ScopeCode string `yaml:"scopeCode"`
	Code      string `yaml:"code"`
}

func (r ConceptReference) String() string {
	return r.ScopeCode + "__" + r.Code
}

type ConceptRepo interface {
	GetConceptByReference(reference ConceptReference) (*Concept, error)
	GetConcept(scopeCode, code string) (*Concept, error)
	ListConceptsForScope(scopeCode string) ([]Concept, error)
}
