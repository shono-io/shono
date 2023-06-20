package kafka

import (
	_ "github.com/benthosdev/benthos/v4/public/components/kafka"
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/utils"
)

func NewInput(opts ...Opt) inventory.Input {
	config := map[string]any{}
	for _, opt := range opts {
		opt(config)
	}

	return inventory.Input{
		Name:       "kafka_franz",
		ConfigSpec: utils.GetBenthosInputConfig("kafka_franz"),
		Config:     config,
	}
}

func NewOutput(opts ...Opt) inventory.Output {
	config := map[string]any{}
	for _, opt := range opts {
		opt(config)
	}

	return inventory.Output{
		Name:       "kafka_franz",
		ConfigSpec: utils.GetBenthosOutputConfig("kafka_franz"),
		Config:     config,
	}
}
