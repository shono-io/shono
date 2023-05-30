package shono

import (
	"fmt"
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func NewBuilder(group string, bb Backbone, reaktors []Reaktor) (*service.StreamBuilder, error) {
	b := service.NewStreamBuilder()

	input, err := toInputYaml(group, bb, reaktors)
	if err != nil {
		return nil, fmt.Errorf("failed to convert input to yaml: %w", err)
	}
	logrus.Debugf("input: %s", input)
	if err := b.AddInputYAML(input); err != nil {
		return nil, fmt.Errorf("failed to add input: %w", err)
	}

	output, err := toOutputYaml(bb, reaktors)
	if err != nil {
		return nil, fmt.Errorf("failed to convert output to yaml: %w", err)
	}
	logrus.Debugf("output: %s", output)
	if err := b.AddOutputYAML(output); err != nil {
		return nil, fmt.Errorf("failed to add output: %w", err)
	}

	processor, err := toProcessorYaml(reaktors)
	if err != nil {
		return nil, fmt.Errorf("failed to convert processor to yaml: %w", err)
	}
	logrus.Debugf("processor: %s", processor)
	if err := b.AddProcessorYAML(processor); err != nil {
		return nil, fmt.Errorf("failed to add processor: %w", err)
	}

	return b, nil
}

func toInputYaml(group string, bb Backbone, reaktors []Reaktor) (string, error) {
	// -- get the list of events from the reaktors
	var events []EventId
	for _, reaktor := range reaktors {
		events = append(events, reaktor.InputEvent())
	}

	res, err := bb.GetConsumerConfig(group, events)
	if err != nil {
		return "", err
	}

	b, err := yaml.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func toOutputYaml(bb Backbone, reaktors []Reaktor) (string, error) {
	// -- get the list of events from the reaktors
	var events []EventId
	for _, reaktor := range reaktors {
		events = append(events, reaktor.OutputEvents()...)
	}

	res, err := bb.GetProducerConfig(events)
	if err != nil {
		return "", err
	}

	b, err := yaml.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func toProcessorYaml(reaktors []Reaktor) (string, error) {
	var cases []map[string]any

	for _, reaktor := range reaktors {
		c, err := toCase(reaktor)
		if err != nil {
			return "", fmt.Errorf("failed to convert reaktor to case: %w", err)
		}

		cases = append(cases, c)
	}

	res := map[string]any{
		"switch": cases,
	}

	b, err := yaml.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func toCase(reaktor Reaktor) (map[string]any, error) {
	res := map[string]any{}

	processor, err := reaktor.Logic().Processor()
	if err != nil {
		return nil, fmt.Errorf("failed to get processor from logic: %w", err)
	}

	// -- we will add a check on the type header
	res["check"] = fmt.Sprintf("@io_shono_kind == %q", reaktor.InputEvent())
	res["processors"] = []map[string]any{processor}

	return res, nil
}
