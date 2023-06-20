package kafka

import (
	"github.com/shono-io/shono/inventory"
)

type SystemType struct {
}

func (s SystemType) Code() string {
	return "kafka"
}

func (s SystemType) RuntimeParameters() []inventory.SystemParameter {
	return []inventory.SystemParameter{
		{
			Name:        "brokers",
			Description: "Comma separated list of Kafka brokers",
		},
	}
}
