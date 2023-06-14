package internal

import (
	"github.com/shono-io/shono/core"
)

type TestSpec struct {
	Summary         string
	EnvironmentVars map[string]any
	Input           core.TestInput
	Assertions      []core.TestAssertion
}

type Test struct {
	Spec TestSpec
}

func (t Test) Summary() string {
	return t.Spec.Summary
}

func (t Test) EnvironmentVars() map[string]any {
	return t.Spec.EnvironmentVars
}

func (t Test) Input() core.TestInput {
	return t.Spec.Input
}

func (t Test) Assertions() []core.TestAssertion {
	return t.Spec.Assertions
}

func HasMetadata(expected map[string]string) TestAssertion {
	return MetadataReaktorTestCondition{
		Values: expected,
		Strict: false,
	}
}

func HasStrictMetadata(expected map[string]string) TestAssertion {
	return MetadataReaktorTestCondition{
		Values: expected,
		Strict: true,
	}
}

type MetadataReaktorTestCondition struct {
	Values map[string]string
	Strict bool
}

func (c MetadataReaktorTestCondition) ConditionType() string {
	return "metadata"
}

func HasPayload(values map[string]interface{}) PayloadReaktorTestCondition {
	return PayloadReaktorTestCondition{
		Values: values,
		Strict: false,
	}
}

func HasStrictPayload(values map[string]interface{}) PayloadReaktorTestCondition {
	return PayloadReaktorTestCondition{
		Values: values,
		Strict: true,
	}
}

type PayloadReaktorTestCondition struct {
	Values map[string]interface{}
	Strict bool
}

func (c PayloadReaktorTestCondition) ConditionType() string {
	return "payload"
}
