package dsl

import (
	"fmt"
	"github.com/shono-io/shono/inventory"
)

func Catch(steps ...inventory.LogicStep) CatchLogicStep {
	return CatchLogicStep{
		Steps: steps,
	}
}

type CatchLogicStep struct {
	Steps []inventory.LogicStep
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
