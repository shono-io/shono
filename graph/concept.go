package graph

import "fmt"

type Concept struct {
	ConceptReference
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Plural      string        `yaml:"plural"`
	Single      string        `yaml:"single"`
	Store       *ConceptStore `yaml:"store"`
	Requests    []Request     `yaml:"requests"`
}

type ConceptStore struct {
	StorageKey string `yaml:"storageKey"`
	Collection string `yaml:"collection"`
}

func ParseConceptReference(input string) (ConceptReference, error) {
	var ref ConceptReference
	if _, err := fmt.Sscanf(input, "%s__%s", &ref.ScopeCode, &ref.Code); err != nil {
		return ConceptReference{}, err
	}

	return ref, nil
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
