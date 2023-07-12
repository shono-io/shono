package dsl

import (
	"fmt"
	"github.com/shono-io/shono/inventory"
)

type TransformBuilder struct {
	*TransformLogicStep
}

func (b *TransformBuilder) Label(label string) *TransformBuilder {
	b.TransformLogicStep.label = label
	return b
}

func (b *TransformBuilder) Mapping(mapping string) *TransformBuilder {
	b.TransformLogicStep.mapping = CleanMultilineStringWhitespace(mapping)
	return b
}

func (b *TransformBuilder) Build() inventory.LogicStep {
	return *b.TransformLogicStep
}

func Transform() *TransformBuilder {
	return &TransformBuilder{&TransformLogicStep{}}
}

type TransformLogicStep struct {
	label   string
	mapping string
}

func (e TransformLogicStep) Validate() error {
	if e.mapping == "" {
		return fmt.Errorf("transform logic must have a mapping")
	}

	return nil
}

func (e TransformLogicStep) Label() string {
	return e.label
}

func (e TransformLogicStep) MarshalBenthos(trace string) (map[string]any, error) {
	trace = fmt.Sprintf("%s/%s", trace, e.Kind())
	if err := e.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", trace, err)
	}

	result := map[string]any{
		"mapping": e.mapping,
	}

	if e.label != "" {
		result["label"] = e.label
	}

	return result, nil
}

func (e TransformLogicStep) Kind() string {
	return "transform"
}
