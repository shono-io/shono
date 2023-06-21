package inventory

import (
	"github.com/shono-io/shono/commons"
)

type Concept interface {
	Node
	Scope() commons.Reference
	Stored() bool
}

type ConceptSpec struct {
	NodeSpec `yaml:",inline"`
	Scope    commons.Reference
	Stored   bool
}

type concept struct {
	Spec ConceptSpec `yaml:",inline"`
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

func (c *concept) Stored() bool {
	return c.Spec.Stored
}
