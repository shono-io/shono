package artifacts

import "github.com/shono-io/shono/inventory"

type GeneratedInput struct {
	Id     string          `json:"id" yaml:"id"`
	Kind   string          `json:"kind" yaml:"kind"`
	Config map[string]any  `json:"config" yaml:"config"`
	Logic  *GeneratedLogic `json:"logic,omitempty" yaml:"logic,omitempty"`
}

type GeneratedOutput struct {
	Id     string         `json:"id" yaml:"id"`
	Kind   string         `json:"kind" yaml:"kind"`
	Config map[string]any `json:"config" yaml:"config"`
}

func AsGeneratedOutput(output inventory.Output) GeneratedOutput {
	return GeneratedOutput{
		Id:     output.Id,
		Kind:   output.Kind,
		Config: output.Config,
	}
}
