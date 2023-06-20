package inventory

type System interface {
	Node
	Spec() SystemSpec
}

type SystemSpec struct {
	Config SystemConfigSpec
}

type SystemConfigSpec struct {
	// Name of the component
	Name string `json:"name"`

	// Type of the component (input, output, etc)
	Type string `json:"type"`

	// A summary of each field in the component configuration.
	Config SystemConfigSpecField `json:"config"`
}

type SystemConfigSpecField struct {
	Name       string                  `json:"name"`
	Type       string                  `json:"type"`
	Kind       string                  `json:"kind"`
	IsOptional bool                    `json:"is_optional,omitempty"`
	Children   []SystemConfigSpecField `json:"children,omitempty"`
}
