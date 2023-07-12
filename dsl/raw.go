package dsl

import (
	"fmt"
	"github.com/shono-io/shono/inventory"
)

type RawBuilder struct {
	*RawLogicStep
}

func (b *RawBuilder) Label(label string) *RawBuilder {
	b.RawLogicStep.label = label
	return b
}

func (b *RawBuilder) Content(content map[string]any) *RawBuilder {
	b.RawLogicStep.Content = content
	return b
}

func (b *RawBuilder) Build() inventory.LogicStep {
	return *b.RawLogicStep
}

func Raw() *RawBuilder {
	return &RawBuilder{&RawLogicStep{}}
}

type RawLogicStep struct {
	label   string
	Content map[string]any
}

func (e RawLogicStep) Label() string {
	return e.label
}

func (e RawLogicStep) MarshalBenthos(trace string) (map[string]any, error) {
	trace = fmt.Sprintf("%s/%s", trace, e.Kind())

	if err := e.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", trace, err)
	}

	result := e.Content

	if e.label != "" {
		result["label"] = e.label
	}

	return result, nil
}

func (e RawLogicStep) Kind() string {
	return "raw"
}

func (e RawLogicStep) Validate() error {
	if e.Content == nil || len(e.Content) == 0 {
		return fmt.Errorf("raw logic must have content")
	}

	return nil
}
