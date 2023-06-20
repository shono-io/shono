package utils

import (
	"encoding/json"
	"fmt"
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/shono-io/shono/inventory"
)

type benthosConfigSpec struct {
	// Name of the component
	Name string `json:"name"`

	// Type of the component (input, output, etc)
	Type string `json:"type"`

	// A summary of each field in the component configuration.
	Config inventory.IOConfigSpecField `json:"config"`
}

func GetBenthosInputConfig(name string) []inventory.IOConfigSpecField {
	var res benthosConfigSpec
	service.GlobalEnvironment().WalkInputs(benthosConfigWalker(name, &res))

	return res.Config.Children
}

func GetBenthosOutputConfig(name string) []inventory.IOConfigSpecField {
	var res benthosConfigSpec
	service.GlobalEnvironment().WalkOutputs(benthosConfigWalker(name, &res))

	return res.Config.Children
}

func benthosConfigWalker(nm string, target any) func(name string, config *service.ConfigView) {
	return func(name string, config *service.ConfigView) {
		if name != nm {
			return
		}

		// -- get the configuration as json
		b, err := config.FormatJSON()
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal(b, &target); err != nil {
			panic(fmt.Errorf("inconsistent benthos expectation: %w", err))
		}
	}
}
