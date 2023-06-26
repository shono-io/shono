package inventory

type TestInput struct {
	Metadata map[string]string
	Content  map[string]any
}

type TestAssertion interface {
	Metadata() map[string]string
	Payload() map[string]any
	Strict() bool
}

type Test struct {
	Summary         string
	EnvironmentVars map[string]any
	Input           TestInput
	Assertions      []TestAssertion
}
