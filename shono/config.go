package shono

type config struct {
	Stream         *StreamConfig         `json:"stream,omitempty"`
	Kafka          *KafkaConfig          `json:"kafka,omitempty"`
	SchemaRegistry *SchemaRegistryConfig `json:"schema_registry,omitempty"`
}

type StreamConfig struct {
	Host  string `json:"host"`
	Token string `json:"token"`
}

type SchemaRegistryConfig struct {
	Urls     []string `json:"urls"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

type KafkaConfig struct {
	Brokers  []string `json:"brokers"`
	Tls      bool     `json:"tls"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}
