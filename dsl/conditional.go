package dsl

import (
	"fmt"
	"github.com/shono-io/shono/inventory"
)

type ConditionalBuilder struct {
	*ConditionalLogicStep
}

func (cb *ConditionalBuilder) Label(label string) *ConditionalBuilder {
	cb.ConditionalLogicStep.label = label
	return cb
}

func (cb *ConditionalBuilder) Case(check string, steps ...inventory.StepBuilder) *ConditionalBuilder {
	cb.ConditionalLogicStep.Cases = append(cb.ConditionalLogicStep.Cases, ConditionalCase{check, inventory.BuildAllSteps(steps...)})

	return cb
}

func (cb *ConditionalBuilder) Default(steps ...inventory.StepBuilder) *ConditionalBuilder {
	cb.ConditionalLogicStep.DefaultCase = &ConditionalCase{Steps: inventory.BuildAllSteps(steps...)}

	return cb
}

func (cb *ConditionalBuilder) Build() inventory.LogicStep {
	return *cb.ConditionalLogicStep
}

func Switch() *ConditionalBuilder {
	return &ConditionalBuilder{&ConditionalLogicStep{}}
}

type ConditionalLogicStep struct {
	label       string
	Cases       []ConditionalCase
	DefaultCase *ConditionalCase
}

func (e ConditionalLogicStep) Label() string {
	return e.label
}

func (e ConditionalLogicStep) MarshalBenthos(trace string) (map[string]any, error) {
	trace = fmt.Sprintf("%s/%s", trace, e.Kind())

	if err := e.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", trace, err)
	}

	var sc []map[string]any

	for cidx, c := range e.Cases {
		// marshall the steps
		var procs []map[string]any
		for sidx, step := range c.Steps {
			t := fmt.Sprintf("%s/cases[%d]/steps[%d]", trace, cidx, sidx)
			proc, err := step.MarshalBenthos(t)
			if err != nil {
				return nil, err
			}
			procs = append(procs, proc)
		}

		sc = append(sc, map[string]any{
			"check":      c.Check,
			"processors": procs,
		})
	}

	if e.DefaultCase != nil {
		// marshall the steps
		var procs []map[string]any
		for sidx, step := range e.DefaultCase.Steps {
			t := fmt.Sprintf("%s/default/steps[%d]", trace, sidx)
			proc, err := step.MarshalBenthos(t)
			if err != nil {
				return nil, err
			}
			procs = append(procs, proc)
		}

		sc = append(sc, map[string]any{
			"processors": procs,
		})
	}

	result := map[string]any{
		"switch": sc,
	}

	if e.label != "" {
		result["label"] = e.label
	}

	return result, nil
}

func (e ConditionalLogicStep) Kind() string {
	return "conditional"
}

func (e ConditionalLogicStep) Validate() error {
	cases := len(e.Cases)
	if e.DefaultCase != nil {
		cases++
	}
	if cases <= 1 {
		return fmt.Errorf("conditional logic must have at least two cases (default included)")
	}

	for _, c := range e.Cases {
		if c.Check == "" {
			return fmt.Errorf("conditional case must have a check")
		}

		if len(c.Steps) == 0 {
			return fmt.Errorf("conditional case must have at least one step")
		}
	}

	if e.DefaultCase != nil {
		if e.DefaultCase.Check != "" {
			return fmt.Errorf("default case cannot have a check")
		}
		if len(e.DefaultCase.Steps) == 0 {
			return fmt.Errorf("default case must have at least one step")
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
