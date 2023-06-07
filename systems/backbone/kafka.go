package backbone

import (
	"github.com/mitchellh/mapstructure"
	"github.com/shono-io/shono/commons"
	"github.com/sirupsen/logrus"
)

type KafkaBackboneConfig struct {
	BootstrapServers []string     `json:"bootstrap_servers" mapstructure:"seed_brokers"`
	TLS              *TLS         `json:"tls" mapstructure:"tls"`
	SASL             []SASLConfig `json:"sasl" mapstructure:"sasl"`

	CheckpointLimit int    `json:"checkpoint_limit" mapstructure:"checkpoint_limit"`
	CommitPeriod    string `json:"commit_period" mapstructure:"commit_period"`
	StartFromOldest bool   `json:"start_from_oldest" mapstructure:"start_from_oldest"`

	LogStrategy LogStrategy `json:"log_strategy" mapstructure:"-"`
}

type TLS struct {
	Enabled             bool      `json:"enabled" mapstructure:"enabled"`
	SkipCertVerify      bool      `json:"skip_cert_verify" mapstructure:"skip_cert_verify"`
	EnableRenegotiation bool      `json:"enable_renegotiation" mapstructure:"enable_renegotiation"`
	RootCas             string    `json:"root_cas" mapstructure:"root_cas"`
	ClientCerts         []TLSCert `json:"client_certs" mapstructure:"client_certs"`
}

type TLSCert struct {
	Cert     string `json:"cert" mapstructure:"cert"`
	Key      string `json:"key" mapstructure:"key"`
	Password string `json:"password" mapstructure:"password"`
}

type AWSConfig struct {
	Region      string         `json:"region" mapstructure:"region"`
	Endpoint    string         `json:"endpoint" mapstructure:"endpoint"`
	Credentials AWSCredentials `json:"credentials" mapstructure:"credentials"`
}

type AWSCredentials struct {
	Profile        string `json:"profile" mapstructure:"profile"`
	Id             string `json:"id" mapstructure:"id"`
	Secret         string `json:"secret" mapstructure:"secret"`
	Token          string `json:"token" mapstructure:"token"`
	FromEC2Role    bool   `json:"from_ec2_role" mapstructure:"from_ec2_role"`
	Role           string `json:"role" mapstructure:"role"`
	RoleExternalId string `json:"role_external_id" mapstructure:"role_external_id"`
}

type SASLConfig struct {
	Mechanism  string            `json:"mechanism" mapstructure:"mechanism"`
	Username   string            `json:"username" mapstructure:"username"`
	Password   string            `json:"password" mapstructure:"password"`
	AwsConfig  *AWSConfig        `json:"aws" mapstructure:"aws"`
	Token      string            `json:"token" mapstructure:"token"`
	Extensions map[string]string `json:"extensions" mapstructure:"extensions"`
}

func NewKafkaBackboneConfig() KafkaBackboneConfig {
	return KafkaBackboneConfig{
		BootstrapServers: []string{"localhost:9092"},
		TLS: &TLS{
			Enabled: false,
		},
		SASL:        []SASLConfig{},
		LogStrategy: PerScopeLogStrategy,
	}
}

func NewKafkaBackbone(cfg KafkaBackboneConfig) Backbone {
	if cfg.LogStrategy == "" {
		cfg.LogStrategy = PerScopeLogStrategy
	}

	return &kafkaBackbone{cfg}
}

type kafkaBackbone struct {
	config KafkaBackboneConfig
}

func (k *kafkaBackbone) GetConsumerConfig(id string, events []commons.Key) (map[string]any, error) {
	var result map[string]any
	if err := mapstructure.Decode(k.config, &result); err != nil {
		return nil, err
	}

	// -- add the topics and the consumer group to the result
	result["topics"] = topicsFromEventIds(events, k.config.LogStrategy)
	result["consumer_group"] = id

	return map[string]any{
		"kafka_franz": result,
	}, nil
}

func (k *kafkaBackbone) GetProducerConfig(events []commons.Key) (map[string]any, error) {
	var result map[string]any
	if err := mapstructure.Decode(k.config, &result); err != nil {
		return nil, err
	}

	var checks []map[string]any
	for _, t := range topicsFromEventIds(events, k.config.LogStrategy) {
		check := map[string]any{
			"check":  "@io_shono_kind.has_prefix(\"" + t + "\")",
			"output": topicOutput(result, t),
		}

		checks = append(checks, check)
	}

	checks = append(checks, map[string]any{
		"output": map[string]any{
			"stdout": map[string]any{
				"codec": "lines",
			},
		},
	})

	return map[string]any{
		"switch": map[string]any{
			"cases": checks,
		},
	}, nil
}

func topicOutput(config map[string]any, topic string) map[string]any {
	result := map[string]any{
		"topic": topic,
	}
	for k, v := range config {
		result[k] = v
	}

	return map[string]any{
		"kafka_franz": result,
	}
}

func topicsFromEventIds(eventIds []commons.Key, strategy LogStrategy) []string {
	topics := map[string]bool{}
	for _, v := range eventIds {
		switch strategy {
		case PerScopeLogStrategy:
			topics[v.Parent().CodeString()] = true
		default:
			logrus.Panicf("unknown log strategy: %v", strategy)
		}
	}

	var result []string
	for k := range topics {
		result = append(result, k)
	}

	return result
}
