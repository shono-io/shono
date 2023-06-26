package inventory

type Input struct {
	Id         string              `json:"id" yaml:"id"`
	Kind       string              `json:"kind" yaml:"kind"`
	ConfigSpec []IOConfigSpecField `json:"-" yaml:"-"`
	Config     map[string]any      `json:"config" yaml:"config"`
	Logic      Logic               `json:"logic,omitempty" yaml:"logic,omitempty"`
}

type Output struct {
	Id         string              `json:"id" yaml:"id"`
	Kind       string              `json:"kind" yaml:"kind"`
	ConfigSpec []IOConfigSpecField `json:"-" yaml:"-"`
	Config     map[string]any      `json:"config" yaml:"config"`
}

type IOConfigSpecField struct {
	Name       string              `json:"name" yaml:"name"`
	Type       string              `json:"type" yaml:"type"`
	Kind       string              `json:"kind,omitempty" yaml:"kind,omitempty"`
	IsOptional bool                `json:"is_optional,omitempty" yaml:"is_optional,omitempty"`
	Children   []IOConfigSpecField `json:"children,omitempty" yaml:"children,omitempty"`
}
