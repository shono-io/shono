package benthos

import (
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/system/kafka"
)

func generateBackboneInput(consumerGroup string, eventRefs []commons.Reference) (*artifacts.GeneratedInput, error) {
	// -- determine the topics to subscribe to
	var topics []string
	fnd := make(map[string]bool)
	for _, ref := range eventRefs {
		t := topicName(ref)
		if !fnd[t] {
			topics = append(topics, t)
			fnd[t] = true
		}
	}

	inp := kafka.NewInput("backbone",
		kafka.WithInputTopics(topics...),
		kafka.WithConsumerGroup(consumerGroup),
	)

	var inpLogic *artifacts.GeneratedLogic
	if inp.Logic != nil {
		l, err := generateLogic(inp.Logic)
		if err != nil {
			return nil, err
		}
		inpLogic = l
	}

	return &artifacts.GeneratedInput{
		Id:     inp.Id,
		Kind:   inp.Kind,
		Config: inp.Config,
		Logic:  inpLogic,
	}, nil
}

func generateBackboneDLQ() (*artifacts.GeneratedOutput, error) {
	out := artifacts.AsGeneratedOutput(kafka.NewOutput(
		"backbone",
		kafka.WithOutputTopic("shono.dlq"),
	))

	return &out, nil
}

func generateBackboneOutput() (*artifacts.GeneratedOutput, error) {
	out := artifacts.AsGeneratedOutput(kafka.NewOutput(
		"backbone",
		kafka.WithOutputTopic("${! @shono_backbone_topic }"),
	))

	return &out, nil
}

func topicName(eventRef commons.Reference) string {
	conceptRef := eventRef.Parent()
	scopeRef := conceptRef.Parent()
	return scopeRef.Code()
}
