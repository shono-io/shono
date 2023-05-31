package benthos

import (
	"fmt"
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/shono-io/shono"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func NewBuilder(group string, bb shono.Backbone, reaktors []shono.Reaktor) (*service.StreamBuilder, error) {
	b := service.NewStreamBuilder()

	input, err := toInputYaml(group, bb, reaktors)
	if err != nil {
		return nil, fmt.Errorf("failed to convert input to yaml: %w", err)
	}
	logrus.Tracef("input: %s", input)
	if err := b.AddInputYAML(input); err != nil {
		return nil, fmt.Errorf("failed to add input: %w", err)
	}

	output, err := toOutputYaml(bb, reaktors)
	if err != nil {
		return nil, fmt.Errorf("failed to convert output to yaml: %w", err)
	}
	logrus.Tracef("output: %s", output)
	if err := b.AddOutputYAML(output); err != nil {
		return nil, fmt.Errorf("failed to add output: %w", err)
	}

	if err := registerCaches(b, reaktors); err != nil {
		return nil, fmt.Errorf("failed to register caches: %w", err)
	}

	processor, err := toProcessorYaml(reaktors)
	if err != nil {
		return nil, fmt.Errorf("failed to convert processor to yaml: %w", err)
	}
	logrus.Tracef("processor: %s", processor)
	if err := b.AddProcessorYAML(processor); err != nil {
		return nil, fmt.Errorf("failed to add processor: %w", err)
	}

	return b, nil
}

func toInputYaml(group string, bb shono.Backbone, reaktors []shono.Reaktor) (string, error) {
	// -- get the list of events from the reaktors
	var events []shono.EventId
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

func toOutputYaml(bb shono.Backbone, reaktors []shono.Reaktor) (string, error) {
	// -- get the list of events from the reaktors
	var events []shono.EventId
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

func toProcessorYaml(reaktors []shono.Reaktor) (string, error) {
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

func toCase(reaktor shono.Reaktor) (map[string]any, error) {
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

func registerCaches(b *service.StreamBuilder, reaktors []shono.Reaktor) error {
	// -- make a list of all stores used by the reaktors
	stores := map[string]shono.Store{}
	for _, reaktor := range reaktors {
		for _, store := range reaktor.Stores() {
			stores[store.Key().String()] = store
		}
	}

	// -- convert each of these stores in their yaml
	for _, store := range stores {
		yml, err := store.AsBenthosComponent()
		if err != nil {
			return fmt.Errorf("failed to convert store %q to yaml: %w", store.Key().String(), err)
		}

		yb, err := yaml.Marshal(yml)
		if err != nil {
			return fmt.Errorf("failed to marshal yaml for store %q: %w", store.Key().String(), err)
		}

		logrus.Tracef("registering cache %s", string(yb))

		if err := b.AddCacheYAML(string(yb)); err != nil {
			return fmt.Errorf("failed to register cache %q: %w", store.Key().String(), err)
		}
	}

	return nil
}
