package artifacts

import "github.com/shono-io/shono/inventory"

type GeneratedInput struct {
	Id         string                        `json:"id" yaml:"id"`
	Kind       string                        `json:"kind" yaml:"kind"`
	ConfigSpec []inventory.IOConfigSpecField `json:"-" yaml:"-"`
	Config     map[string]any                `json:"config" yaml:"config"`
	Logic      *GeneratedLogic               `json:"logic,omitempty" yaml:"logic,omitempty"`
}
