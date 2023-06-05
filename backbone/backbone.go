package backbone

import (
	"github.com/shono-io/shono/commons"
	"github.com/sirupsen/logrus"
)

type LogStrategy string

var (
	PerScopeLogStrategy LogStrategy = "per_scope"
)

type Backbone interface {
	GetConsumerConfig(id string, events []commons.Key) (map[string]any, error)
	GetProducerConfig(events []commons.Key) (map[string]any, error)
}

func NewKafkaBackbone(config map[string]any, logStrategy LogStrategy) Backbone {
	return &kafkaBackbone{
		config:      config,
		logStrategy: logStrategy,
	}
}

type kafkaBackbone struct {
	config      map[string]any
	logStrategy LogStrategy
}

func (k *kafkaBackbone) GetConsumerConfig(id string, events []commons.Key) (map[string]any, error) {
	result := map[string]any{
		"topics":         topicsFromEventIds(events, k.logStrategy),
		"consumer_group": id,
	}
	for k, v := range k.config {
		result[k] = v
	}

	return map[string]any{
		"kafka_franz": result,
	}, nil
}

func (k *kafkaBackbone) GetProducerConfig(events []commons.Key) (map[string]any, error) {
	var checks []map[string]any
	for _, t := range topicsFromEventIds(events, k.logStrategy) {
		check := map[string]any{
			"check":  "@io_shono_kind.has_prefix(\"" + t + "\")",
			"output": topicOutput(k.config, t),
		}

		checks = append(checks, check)
	}

	checks = append(checks, map[string]any{
		"output": map[string]any{
			"stdout": map[string]any{
				"codec": "lines",
			},
		},
	})

	return map[string]any{
		"switch": map[string]any{
			"cases": checks,
		},
	}, nil
}

func topicOutput(config map[string]any, topic string) map[string]any {
	result := map[string]any{
		"topic": topic,
	}
	for k, v := range config {
		result[k] = v
	}

	return map[string]any{
		"kafka_franz": result,
	}
}

func topicsFromEventIds(eventIds []commons.Key, strategy LogStrategy) []string {
	topics := map[string]bool{}
	for _, v := range eventIds {
		switch strategy {
		case PerScopeLogStrategy:
			topics[v.Parent().CodeString()] = true
		default:
			logrus.Panicf("unknown log strategy: %v", strategy)
		}
	}

	var result []string
	for k := range topics {
		result = append(result, k)
	}

	return result
}
