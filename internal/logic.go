package internal

import "github.com/shono-io/shono/core"

type LogicSpec struct {
	Steps []core.LogicStep
	Tests []core.Test
}

type Logic struct {
	Spec LogicSpec
}

func (l Logic) Steps() []core.LogicStep {
	return l.Spec.Steps
}

func (l Logic) Tests() []core.Test {
	return l.Spec.Tests
}
