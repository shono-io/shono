package dsl

import (
	"fmt"
	"github.com/shono-io/shono/core"
)

func Catch(steps ...core.LogicStep) CatchLogicStep {
	return CatchLogicStep{
		Steps: steps,
	}
}

type CatchLogicStep struct {
	Steps []core.LogicStep
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
