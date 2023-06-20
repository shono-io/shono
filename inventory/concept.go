package inventory

import (
	"github.com/shono-io/shono/commons"
)

type Concept interface {
	Node
	Scope() commons.Reference
	Store() *ConceptStore
}

type ConceptStore struct {
	Storage    Storage
	Collection string
}

type ConceptSpec struct {
	NodeSpec
	Scope commons.Reference
	Store *ConceptStore
}

type concept struct {
	Spec ConceptSpec
}

func (c *concept) Code() string {
	return c.Spec.Code
}

func (c *concept) Summary() string {
	return c.Spec.Summary
}

func (c *concept) Docs() string {
	return c.Spec.Docs
}

func (c *concept) Status() commons.Status {
	return c.Spec.Status
}

func (c *concept) Scope() commons.Reference {
	return c.Spec.Scope
}

func (c *concept) Reference() commons.Reference {
	return c.Scope().Child("concepts", c.Code())
}

func (c *concept) Store() *ConceptStore {
	return c.Spec.Store
}
