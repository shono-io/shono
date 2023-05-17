package shono

type Config struct {
	Stream   *StreamConfig   `json:"stream,omitempty"`
	Backbone *BackboneConfig `json:"backbone,omitempty"`
}

type StreamConfig struct {
	Host  string `json:"host"`
	Token string `json:"token"`
}

type BackboneConfig struct {
	Kind       string         `json:"kind"`
	Properties map[string]any `json:"properties"`
}
