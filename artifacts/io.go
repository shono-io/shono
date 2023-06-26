package artifacts

import "github.com/shono-io/shono/inventory"

type GeneratedInput struct {
	Name       string                        `json:"name" yaml:"name"`
	ConfigSpec []inventory.IOConfigSpecField `json:"-" yaml:"-"`
	Config     map[string]any                `json:"config" yaml:"config"`
	Logic      *GeneratedLogic               `json:"logic,omitempty" yaml:"logic,omitempty"`
}
