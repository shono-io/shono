package inventory

func NewTest(summary string) *TestBuilder {
	return &TestBuilder{
		test: &Test{
			Summary: summary,
		},
	}
}

type TestBuilder struct {
	test *Test
}

func (b *TestBuilder) Given(environment map[string]any) *TestBuilder {
	b.test.EnvironmentVars = environment
	return b
}

func (b *TestBuilder) When(input TestInput) *TestBuilder {
	b.test.Input = input
	return b
}

func (b *TestBuilder) Then(assertions ...TestAssertion) *TestBuilder {
	b.test.Assertions = append(b.test.Assertions, assertions...)
	return b
}

func (b *TestBuilder) Build() Test {
	return *b.test
}
