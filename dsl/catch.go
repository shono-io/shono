package dsl

import (
	"fmt"
	"github.com/shono-io/shono/inventory"
)

type CatchOpt func(c *CatchLogicStep)

func WithCatchLabel(label string) CatchOpt {
	return func(c *CatchLogicStep) {
		c.label = label
	}
}

type CatchBuilder struct {
	result *CatchLogicStep
}

func (cb *CatchBuilder) Label(label string) *CatchBuilder {
	cb.result.label = label
	return cb
}

func (cb *CatchBuilder) Steps(steps ...inventory.StepBuilder) *CatchBuilder {
	cb.result.Steps = append(cb.result.Steps, inventory.BuildAllSteps(steps...)...)
	return cb
}

func (cb *CatchBuilder) Build() inventory.LogicStep {
	return *cb.result
}

func Catch() *CatchBuilder {
	return &CatchBuilder{
		result: &CatchLogicStep{},
	}
}

type CatchLogicStep struct {
	label string
	Steps []inventory.LogicStep
}

func (e CatchLogicStep) Label() string {
	return e.label
}

func (e CatchLogicStep) MarshalBenthos(trace string) (map[string]any, error) {
	trace = fmt.Sprintf("%s/%s", trace, e.Kind())

	if err := e.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", trace, err)
	}

	// marshall the steps
	var procs []map[string]any
	for idx, step := range e.Steps {
		t := fmt.Sprintf("%s/clauses[%d]", trace, idx)
		proc, err := step.MarshalBenthos(t)
		if err != nil {
			return nil, err
		}
		procs = append(procs, proc)
	}

	result := map[string]any{
		"catch": procs,
	}

	if e.label != "" {
		result["label"] = e.label
	}

	return result, nil
}

func (e CatchLogicStep) Kind() string {
	return "catch"
}

func (e CatchLogicStep) Validate() error {
	if len(e.Steps) == 0 {
		return fmt.Errorf("catch logic must have at least one step")
	}

	return nil
}
