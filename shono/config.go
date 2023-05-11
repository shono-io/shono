package shono

type config struct {
	Kafka          KafkaConfig          `json:"kafka"`
	SchemaRegistry SchemaRegistryConfig `json:"schema_registry"`
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
