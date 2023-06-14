package kafka

type Config struct {
	BootstrapServers []string     `yaml:"bootstrap_servers" mapstructure:"seed_brokers"`
	TLS              *TLS         `yaml:"tls,omitempty" mapstructure:"tls,omitempty"`
	SASL             []SASLConfig `yaml:"sasl,omitempty" mapstructure:"sasl,omitempty"`

	CheckpointLimit *int    `yaml:"checkpoint_limit,omitempty" mapstructure:"checkpoint_limit,omitempty" yaml:"checkpoint_limit,omitempty"`
	CommitPeriod    *string `yaml:"commit_period,omitempty" mapstructure:"commit_period,omitempty" yaml:"commit_period,omitempty"`
	StartFromOldest *bool   `yaml:"start_from_oldest,omitempty" mapstructure:"start_from_oldest,omitempty" yaml:"start_from_oldest,omitempty"`
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
