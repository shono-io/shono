package file

import (
	_ "github.com/benthosdev/benthos/v4/public/components/pure"
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/utils"
)

func NewInput(opts ...Opt) inventory.Input {
	config := map[string]any{}
	for _, opt := range opts {
		opt(config)
	}

	return inventory.Input{
		Name:       "file",
		ConfigSpec: utils.GetBenthosInputConfig("file"),
		Config:     config,
	}
}

func NewOutput(opts ...Opt) inventory.Output {
	config := map[string]any{}
	for _, opt := range opts {
		opt(config)
	}

	return inventory.Output{
		Name:       "file",
		ConfigSpec: utils.GetBenthosOutputConfig("file"),
		Config:     config,
	}
}
