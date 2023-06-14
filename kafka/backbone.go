package kafka

import (
	"github.com/mitchellh/mapstructure"
	"github.com/shono-io/shono/core"
)

type Backbone struct {
	Config Config
}

func (b Backbone) EventLogName(event core.Event) string {
	// -- currently only supporting event logs on scope level
	scopeReference := event.Concept().Parent()
	return scopeReference.Code()
}

func (b Backbone) AsInput(id string, events ...core.Event) (map[string]any, error) {
	var result map[string]any
	if err := mapstructure.Decode(b.Config, &result); err != nil {
		return nil, err
	}

	// -- add the topics and the consumer group to the result
	result["topics"] = b.topicsFromEventIds(events)
	result["consumer_group"] = id

	return map[string]any{
		"kafka_franz": result,
	}, nil
}

func (b Backbone) AsOutput(events ...core.Event) (map[string]any, error) {
	var result map[string]any
	if err := mapstructure.Decode(b.Config, &result); err != nil {
		return nil, err
	}

	var checks []map[string]any
	for _, t := range b.topicsFromEventIds(events) {
		check := map[string]any{
			"check":  "@io_shono_kind.has_prefix(\"" + t + "\")",
			"output": topicOutput(result, t),
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

	delete(result, "checkpoint_limit")
	delete(result, "commit_period")
	delete(result, "start_from_oldest")

	return map[string]any{
		"kafka_franz": result,
	}
}

func (b Backbone) topicsFromEventIds(events []core.Event) []string {
	var result []string
	fnd := map[string]bool{}

	for _, v := range events {
		t := b.EventLogName(v)
		if _, ok := fnd[t]; !ok {
			result = append(result, t)
			fnd[t] = true
		}
	}

	return result
}
