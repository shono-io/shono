package core

import "github.com/shono-io/shono/internal"

func NewTest(summary string) TestBuilder {
	return TestBuilder{
		spec: internal.TestSpec{
			Summary: summary,
		},
	}
}

type TestBuilder struct {
	spec internal.TestSpec
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
	return internal.Test{b.spec}
}

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

func EventInput(eventRef Reference, content map[string]any) TestInput {
	return TestInput{
		Metadata: map[string]string{
			"io_shono_kind": eventRef.String(),
		},
		Content: content,
	}
}

func RawInput(metadata map[string]string, content map[string]any) TestInput {
	return TestInput{
		Metadata: metadata,
		Content:  content,
	}
}

type TestAssertion interface {
	ConditionType() string
}

func AssertMetadataEquals(expected map[string]string) TestAssertion {
	return MetadataReaktorTestCondition{
		Values: expected,
		Strict: false,
	}
}

func AssertMetadataContains(key, value string) TestAssertion {
	return MetadataReaktorTestCondition{
		Values: map[string]string{
			key: value,
		},
		Strict: false,
	}
}

func AssertContentEquals(expected map[string]interface{}) TestAssertion {
	return PayloadReaktorTestCondition{
		Values: expected,
	}
}

func AssertContentContains(key string, value interface{}) TestAssertion {
	return PayloadReaktorTestCondition{
		Values: map[string]interface{}{
			key: value,
		},
	}
}
