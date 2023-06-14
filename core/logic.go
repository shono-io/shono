package core

import "github.com/shono-io/shono/internal"

func NewLogic() LogicBuilder {
	return LogicBuilder{}
}

type LogicBuilder struct {
	spec internal.LogicSpec
}

func (b LogicBuilder) Steps(steps ...LogicStep) LogicBuilder {
	b.spec.Steps = steps
	return b
}

func (b LogicBuilder) Build() Logic {
	return internal.Logic{b.spec}
}

type Logic interface {
	Steps() []LogicStep
	Tests() []Test
}

type LogicStep interface {
	Kind() string
	Validate() error
}
