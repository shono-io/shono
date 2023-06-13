package backbone

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/shono-io/shono/graph"
	"github.com/sirupsen/logrus"
)

type KafkaBackboneConfig struct {
	BootstrapServers []string     `yaml:"bootstrap_servers" mapstructure:"seed_brokers"`
	TLS              *TLS         `yaml:"tls,omitempty" mapstructure:"tls,omitempty"`
	SASL             []SASLConfig `yaml:"sasl,omitempty" mapstructure:"sasl,omitempty"`

	CheckpointLimit *int    `yaml:"checkpoint_limit,omitempty" mapstructure:"checkpoint_limit,omitempty" yaml:"checkpoint_limit,omitempty"`
	CommitPeriod    *string `yaml:"commit_period,omitempty" mapstructure:"commit_period,omitempty" yaml:"commit_period,omitempty"`
	StartFromOldest *bool   `yaml:"start_from_oldest,omitempty" mapstructure:"start_from_oldest,omitempty" yaml:"start_from_oldest,omitempty"`

	LogStrategy graph.LogStrategy `yaml:"log_strategy" mapstructure:"-"`
}

type TLS struct {
	Enabled             bool      `yaml:"enabled,omitempty" mapstructure:"enabled,omitempty"`
	SkipCertVerify      bool      `yaml:"skip_cert_verify,omitempty" mapstructure:"skip_cert_verify,omitempty"`
	EnableRenegotiation bool      `yaml:"enable_renegotiation,omitempty" mapstructure:"enable_renegotiation,omitempty"`
	RootCas             string    `yaml:"root_cas,omitempty" mapstructure:"root_cas,omitempty"`
	ClientCerts         []TLSCert `yaml:"client_certs,omitempty" mapstructure:"client_certs,omitempty"`
}

type TLSCert struct {
	Cert     string `yaml:"cert,omitempty" mapstructure:"cert,omitempty"`
	Key      string `yaml:"key,omitempty" mapstructure:"key,omitempty"`
	Password string `yaml:"password,omitempty" mapstructure:"password,omitempty"`
}

type AWSConfig struct {
	Region      string         `yaml:"region,omitempty" mapstructure:"region,omitempty"`
	Endpoint    string         `yaml:"endpoint,omitempty" mapstructure:"endpoint,omitempty"`
	Credentials AWSCredentials `yaml:"credentials,omitempty" mapstructure:"credentials,omitempty"`
}

type AWSCredentials struct {
	Profile        string `yaml:"profile,omitempty" mapstructure:"profile,omitempty"`
	Id             string `yaml:"id,omitempty" mapstructure:"id,omitempty"`
	Secret         string `yaml:"secret,omitempty" mapstructure:"secret,omitempty"`
	Token          string `yaml:"token,omitempty" mapstructure:"token,omitempty"`
	FromEC2Role    bool   `yaml:"from_ec2_role,omitempty" mapstructure:"from_ec2_role,omitempty"`
	Role           string `yaml:"role,omitempty" mapstructure:"role,omitempty"`
	RoleExternalId string `yaml:"role_external_id,omitempty" mapstructure:"role_external_id,omitempty"`
}

type SASLConfig struct {
	Mechanism  string            `yaml:"mechanism,omitempty" mapstructure:"mechanism,omitempty"`
	Username   string            `yaml:"username,omitempty" mapstructure:"username,omitempty"`
	Password   string            `yaml:"password,omitempty" mapstructure:"password,omitempty"`
	Aws        AWSConfig         `yaml:"aws,omitempty" mapstructure:"aws,omitempty"`
	Token      string            `yaml:"token,omitempty" mapstructure:"token,omitempty"`
	Extensions map[string]string `yaml:"extensions,omitempty" mapstructure:"extensions,omitempty"`
}

type kafkaBackbone struct {
	config KafkaBackboneConfig
}

func (k *kafkaBackbone) GetConsumerConfig(id string, events []graph.EventReference) (map[string]any, error) {
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

func (k *kafkaBackbone) GetProducerConfig(events []graph.EventReference) (map[string]any, error) {
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

	delete(result, "checkpoint_limit")
	delete(result, "commit_period")
	delete(result, "start_from_oldest")

	return map[string]any{
		"kafka_franz": result,
	}
}

func topicsFromEventIds(eventIds []graph.EventReference, strategy graph.LogStrategy) []string {
	topics := map[string]bool{}
	for _, v := range eventIds {
		switch strategy {
		case graph.PerScopeLogStrategy:
			topics[fmt.Sprintf("%s__%s", v.ScopeCode, v.ConceptCode)] = true
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
