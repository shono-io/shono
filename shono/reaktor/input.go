package reaktor

import (
	"github.com/shono-io/go-shono/shono"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func NewShonoInput(id string, backbone shono.Backbone, eventIds ...shono.EventId) *Input {
	return &Input{
		Id:       id,
		EventIds: eventIds,
		Backbone: backbone,
	}
}

type Input struct {
	Id       string
	EventIds []shono.EventId
	Backbone shono.Backbone
}

func (s *Input) AsYaml() string {
	v, err := s.Backbone.GetConsumerConfig(s.Id, s.EventIds)
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
