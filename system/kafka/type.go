package kafka

import (
	_ "github.com/benthosdev/benthos/v4/public/components/kafka"
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/utils"
)

func NewInput(id string, opts ...Opt) inventory.Input {
	config := map[string]any{}

	for _, opt := range opts {
		opt(config)
	}

	return inventory.Input{
		Id:         id,
		Kind:       "kafka_franz",
		ConfigSpec: utils.GetBenthosInputConfig("kafka_franz"),
		Config:     config,
		Logic: inventory.NewLogic().Steps(
			dsl.Log("INFO", "RAW metadata ${!@} with payload ${! json(\"key\") }"),
			dsl.Transform(dsl.BloblangMapping(`
				root = this
				meta shono_key = @kafka_key
				meta shono_timestamp = @timestamp_unix
				meta = @.filter(kv -> !kv.key.has_prefix("kafka_"))
			`)),
			dsl.Log("INFO", "RAW postprocessed metadata ${!@} with payload ${! json(\"key\") }"),
		).Build(),
	}
}

func NewOutput(id string, opts ...Opt) inventory.Output {
	config := map[string]any{
		"key": "${! @shono_key }",
		"metadata": map[string]any{
			"include_prefixes": []string{"shono_"},
		},
	}
	for _, opt := range opts {
		opt(config)
	}

	return inventory.Output{
		Id:         id,
		Kind:       "kafka_franz",
		ConfigSpec: utils.GetBenthosOutputConfig("kafka_franz"),
		Config:     config,
	}
}
