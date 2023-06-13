package backbone

import (
	"fmt"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
)

type Config struct {
	Kind   string                `yaml:"kind"`
	Config commons.PartialConfig `yaml:"config"`
}

func NewBackbone(config Config) (graph.Backbone, error) {
	switch config.Kind {
	case "kafka":
		var cfg KafkaBackboneConfig
		if err := config.Config.Unmarshal(&cfg); err != nil {
			return nil, err
		}

		if cfg.LogStrategy == "" {
			cfg.LogStrategy = graph.PerScopeLogStrategy
		}

		return &kafkaBackbone{cfg}, nil

	default:
		return nil, fmt.Errorf("unknown backbone kind: %s", config.Kind)
	}
}
