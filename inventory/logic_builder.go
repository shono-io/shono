package inventory

func NewLogic() *LogicBuilder {
	return &LogicBuilder{}
}

type LogicBuilder struct {
	spec LogicSpec
}

func (b *LogicBuilder) Steps(steps ...StepBuilder) *LogicBuilder {
	b.spec.Steps = steps
	return b
}

func (b *LogicBuilder) Test(tests ...Test) *LogicBuilder {
	b.spec.Tests = append(b.spec.Tests, tests...)
	return b
}

func (b *LogicBuilder) Build() Logic {
	return logic{b.spec}
}
