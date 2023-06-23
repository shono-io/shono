package inventory

import (
	"github.com/shono-io/shono/commons"
)

type Inventory interface {
	ResolveScope(ref commons.Reference) (Scope, error)
	ResolveConcept(ref commons.Reference) (Concept, error)
	ResolveEvent(ref commons.Reference) (Event, error)
	ResolveInjector(ref commons.Reference) (Injector, error)
	ResolveExtractor(ref commons.Reference) (Extractor, error)

	ListInjectorsForScope(scopeRef commons.Reference) ([]Injector, error)
	ListReactorsForConcept(conceptRef commons.Reference) ([]Reactor, error)
	ListExtractorsForScope(scopeRef commons.Reference) ([]Extractor, error)
}

type Node interface {
	Code() string
	Summary() string
	Docs() string
	Status() commons.Status
	Reference() commons.Reference
}

type Executable interface {
	Node
	Logic() Logic
}

type NodeSpec struct {
	Code    string
	Summary string
	Docs    string
	Status  commons.Status
}

type node struct {
	spec NodeSpec
}

func (n *node) Code() string {
	return n.spec.Code
}

func (n *node) Summary() string {
	return n.spec.Summary
}

func (n *node) Docs() string {
	return n.spec.Docs
}

func (n *node) Status() commons.Status {
	return n.spec.Status
}
