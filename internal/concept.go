package internal

import (
	"github.com/shono-io/shono/core"
)

type ConceptSpec struct {
	NodeSpec
	Scope core.Reference
}

type Concept struct {
	Spec ConceptSpec
}

func (c *Concept) Code() string {
	return c.Spec.Code
}

func (c *Concept) Summary() string {
	return c.Spec.Summary
}

func (c *Concept) Docs() string {
	return c.Spec.Docs
}

func (c *Concept) Status() core.Status {
	return c.Spec.Status
}

func (c *Concept) Scope() core.Reference {
	return c.Spec.Scope
}

func (c *Concept) Reference() core.Reference {
	return c.Scope().Child("concepts", c.Code())
}
