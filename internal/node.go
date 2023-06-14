package internal

import (
	"github.com/shono-io/shono/core"
)

type NodeSpec struct {
	Code    string
	Summary string
	Docs    string
	Status  core.Status
}

type Node struct {
	spec NodeSpec
}

func (n *Node) Code() string {
	return n.spec.Code
}

func (n *Node) Summary() string {
	return n.spec.Summary
}

func (n *Node) Docs() string {
	return n.spec.Docs
}

func (n *Node) Status() core.Status {
	return n.spec.Status
}
