package inventory

type Logic interface {
	Steps() []LogicStep
	Tests() []Test
}

type LogicStep interface {
	Kind() string
	Validate() error
}

type LogicSpec struct {
	Steps []LogicStep
	Tests []Test
}

type logic struct {
	Spec LogicSpec `yaml:",inline"`
}

func (l logic) Steps() []LogicStep {
	return l.Spec.Steps
}

func (l logic) Tests() []Test {
	return l.Spec.Tests
}
