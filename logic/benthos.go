package logic

import (
	"gopkg.in/yaml.v3"
)

var Benthos EngineId = "benthos"

func NewBenthosLogic(source string) Logic {
	return &benthosLogic{source: source}
}

type benthosLogic struct {
	source string
}

func (b *benthosLogic) EngineId() EngineId {
	return Benthos
}

func (b *benthosLogic) Processor() (map[string]any, error) {
	var result map[string]any

	// -- parse the source as being a benthos processor config
	if err := yaml.Unmarshal([]byte(b.source), &result); err != nil {
		return nil, err
	}

	return result, nil
}
