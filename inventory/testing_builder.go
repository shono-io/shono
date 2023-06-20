package inventory

func NewTest(summary string) TestBuilder {
	return TestBuilder{
		spec: TestSpec{
			Summary: summary,
		},
	}
}

type TestBuilder struct {
	spec TestSpec
}

func (b TestBuilder) Given(environment map[string]any) TestBuilder {
	b.spec.EnvironmentVars = environment
	return b
}

func (b TestBuilder) When(input TestInput) TestBuilder {
	b.spec.Input = input
	return b
}

func (b TestBuilder) Then(assertions ...TestAssertion) TestBuilder {
	b.spec.Assertions = append(b.spec.Assertions, assertions...)
	return b
}

func (b TestBuilder) Build() Test {
	return test{b.spec}
}
