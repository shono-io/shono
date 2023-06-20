package inventory

type Test interface {
	Summary() string
	EnvironmentVars() map[string]any
	Input() TestInput
	Assertions() []TestAssertion
}

type TestInput struct {
	Metadata map[string]string
	Content  map[string]any
}

type TestAssertion interface {
	Metadata() map[string]string
	Payload() map[string]any
	Strict() bool
}

type TestSpec struct {
	Summary         string
	EnvironmentVars map[string]any
	Input           TestInput
	Assertions      []TestAssertion
}

type test struct {
	Spec TestSpec
}

func (t test) Summary() string {
	return t.Spec.Summary
}

func (t test) EnvironmentVars() map[string]any {
	return t.Spec.EnvironmentVars
}

func (t test) Input() TestInput {
	return t.Spec.Input
}

func (t test) Assertions() []TestAssertion {
	return t.Spec.Assertions
}
