package dsl

import (
	"fmt"
	"github.com/shono-io/shono/inventory"
)

func Switch(cases ...ConditionalCase) ConditionalLogicStep {
	return ConditionalLogicStep{
		Cases: cases,
	}
}

func SwitchCase(check string, steps ...inventory.LogicStep) ConditionalCase {
	return ConditionalCase{
		Check: check,
		Steps: steps,
	}
}

func SwitchDefault(steps ...inventory.LogicStep) ConditionalCase {
	return ConditionalCase{
		Steps: steps,
	}
}

type ConditionalLogicStep struct {
	Cases []ConditionalCase `yaml:"cases"`
}

func (e ConditionalLogicStep) Kind() string {
	return "conditional"
}

func (e ConditionalLogicStep) Validate() error {
	if len(e.Cases) <= 1 {
		return fmt.Errorf("conditional logic must have at least two cases")
	}

	for _, c := range e.Cases {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type ConditionalCase struct {
	Check string                `yaml:"check,omitempty"`
	Steps []inventory.LogicStep `yaml:"steps"`
}

func (e ConditionalCase) Validate() error {
	if len(e.Steps) == 0 {
		return fmt.Errorf("conditional case must have at least one step")
	}

	return nil
}
