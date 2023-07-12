package inventory

type Logic interface {
	Steps() []StepBuilder
	Tests() []Test
}

type LogicStep interface {
	Label() string
	Kind() string
	Validate() error
	MarshalBenthos(trace string) (map[string]any, error)
}

type LogicSpec struct {
	Steps []StepBuilder
	Tests []Test
}

type StepBuilder interface {
	Build() LogicStep
}

func BuildAllSteps(steps ...StepBuilder) []LogicStep {
	result := make([]LogicStep, len(steps))
	for i, step := range steps {
		result[i] = step.Build()
	}
	return result
}

type logic struct {
	Spec LogicSpec `yaml:",inline"`
}

func (l logic) Steps() []StepBuilder {
	return l.Spec.Steps
}

func (l logic) Tests() []Test {
	return l.Spec.Tests
}
