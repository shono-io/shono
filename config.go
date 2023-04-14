package go_shono

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type AppConfig struct {
	Domain string      `json:"domain" yaml:"domain"`
	AppId  string      `json:"app_id" yaml:"app_id"`
	Kafka  KafkaConfig `json:"kafka" yaml:"kafka"`
}

type KafkaConfig struct {
	Brokers []string `json:"brokers" yaml:"brokers"`
}

func (c AppConfig) KafkaOpts() []kgo.Opt {
	return []kgo.Opt{
		kgo.SeedBrokers(c.Kafka.Brokers...),
		kgo.ConsumerGroup(c.AppId),
	}
}
