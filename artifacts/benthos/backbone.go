package benthos

import (
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/system/kafka"
)

func generateBackboneInput(eventRefs []commons.Reference) (*artifacts.GeneratedInput, error) {
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

	inp := kafka.NewInput(
		kafka.WithInputTopics(topics...),
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
		Name:       inp.Name,
		ConfigSpec: inp.ConfigSpec,
		Config:     inp.Config,
		Logic:      inpLogic,
	}, nil
}

func generateBackboneDLQ() (inventory.Output, error) {
	out := kafka.NewOutput(
		kafka.WithOutputTopic("shono.dlq"),
	)
	return out, nil
}

func generateBackboneOutput() (inventory.Output, error) {
	out := kafka.NewOutput(
		kafka.WithOutputTopic("${!@backbone_topic}"),
	)
	return out, nil
}

func topicName(eventRef commons.Reference) string {
	conceptRef := eventRef.Parent()
	scopeRef := conceptRef.Parent()
	return scopeRef.Code()
}
