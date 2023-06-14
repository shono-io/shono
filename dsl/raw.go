package dsl

import "fmt"

func Raw(content map[string]any) RawLogicStep {
	return RawLogicStep{
		Content: content,
	}
}

type RawLogicStep struct {
	Content map[string]any `yaml:"content"`
}

func (e RawLogicStep) Kind() string {
	return "raw"
}

func (e RawLogicStep) Validate() error {
	if e.Content == nil {
		return fmt.Errorf("raw logic must have content")
	}

	return nil
}
