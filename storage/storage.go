package storage

import (
	"fmt"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
)

type Config struct {
	Kind   string                `yaml:"kind"`
	Config commons.PartialConfig `yaml:"config"`
}

func NewStorage(key string, config Config) (graph.Storage, error) {
	switch config.Kind {
	case "arangodb":
		var cfg ArangodbConfig
		if err := config.Config.Unmarshal(&cfg); err != nil {
			return nil, err
		}

		return &arangodb{key, cfg}, nil
	case "mongodb":
		var cfg MongodbConfig
		if err := config.Config.Unmarshal(&cfg); err != nil {
			return nil, err
		}

		return &mongodb{key, cfg}, nil
	default:
		return nil, fmt.Errorf("unknown storage system kind: %s", config.Kind)
	}
}
