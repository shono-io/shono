package reaktor

import (
	"github.com/shono-io/go-shono/shono"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func NewOutput(backbone shono.Backbone, eventIds ...shono.EventId) *Output {
	return &Output{
		EventIds: eventIds,
		Backbone: backbone,
	}
}

type Output struct {
	EventIds []shono.EventId
	Backbone shono.Backbone
}

func (o *Output) AsYaml() string {
	v, err := o.Backbone.GetProducerConfig(o.EventIds)
	if err != nil {
		logrus.Panicf("failed to get consumer config: %v", err)
	}

	// -- marshal
	b, err := yaml.Marshal(v)
	if err != nil {
		logrus.Panicf("failed to marshal inbound reaktor: %v", err)
	}

	return string(b)
}
