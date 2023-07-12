package benthos

import (
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/utils"
)

func NewInput(id string, kind string, config map[string]any) *inventory.Input {
	return &inventory.Input{
		Id:         id,
		Kind:       kind,
		ConfigSpec: utils.GetBenthosInputConfig(kind),
		Config:     config,
	}
}

func NewOutput(id string, kind string, config map[string]any) *inventory.Output {
	return &inventory.Output{
		Id:         id,
		Kind:       kind,
		ConfigSpec: utils.GetBenthosOutputConfig(kind),
		Config:     config,
	}
}
