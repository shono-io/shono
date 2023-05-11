package shono

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sr"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type ClientConfig struct {
	Org    string `json:"org"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Url    string `json:"url"`
}

type ClientOpt func(*ClientConfig)

func WithOrg(org string) ClientOpt {
	return func(cfg *ClientConfig) {
		cfg.Org = org
	}
}

func WithCredentials(key string, secret string) ClientOpt {
	return func(cfg *ClientConfig) {
		cfg.Key = key
		cfg.Secret = secret
	}
}

func WithUrl(url string) ClientOpt {
	return func(cfg *ClientConfig) {
		cfg.Url = url
	}
}

func NewClient(opts ...ClientOpt) (*Client, error) {
	co := &ClientConfig{
		Url: "https://dev-api.shono.io",
	}

	for _, opt := range opts {
		opt(co)
	}

	at, err := getToken(co.Key, co.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	cfg, err := getConfig(at, co.Org, co.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to get config from shono: %w", err)
	}

	kc, err := getKafkaClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %w", err)
	}

	src, err := getSchemaRegistryClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create schema registry client: %w", err)
	}

	return &Client{kc, src}, nil
}

type Client struct {
	kc  *kgo.Client
	src *sr.Client
}

func (c *Client) Kafka() *kgo.Client {
	return c.kc
}

func (c *Client) SchemaRegistry() *sr.Client {
	return c.src
}

func (c *Client) Close() {
	c.kc.Close()
}

func getKafkaClient(config *config) (*kgo.Client, error) {
	var opts []kgo.Opt
	opts = append(opts, kgo.SeedBrokers(config.Kafka.Brokers...))

	if config.Kafka.Username != "" && config.Kafka.Password != "" {
		opts = append(opts, kgo.SASL(plain.Auth{User: config.Kafka.Username, Pass: config.Kafka.Password}.AsMechanism()))
	}

	if config.Kafka.Tls {
		tlsDialer := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}
		opts = append(opts, kgo.Dialer(tlsDialer.DialContext))
	}

	return kgo.NewClient(opts...)
}

func getSchemaRegistryClient(config *config) (*sr.Client, error) {
	var opts []sr.Opt
	opts = append(opts, sr.URLs(config.SchemaRegistry.Urls...))

	if config.SchemaRegistry.Username != "" && config.SchemaRegistry.Password != "" {
		opts = append(opts, sr.BasicAuth(config.SchemaRegistry.Username, config.SchemaRegistry.Password))
	}

	return sr.NewClient(opts...)
}

func getConfig(at string, org string, baseUrl string) (*config, error) {
	url := fmt.Sprintf("%s/organizations/%s/config", baseUrl, org)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", at))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get config from shono: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to get config from shono: %s", res.Status)
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	var body config
	if err := json.Unmarshal(b, &body); err != nil {
		return nil, err
	}

	return &body, nil
}

func getToken(shonoKey string, shonoSecret string) (string, error) {
	url := "https://dev-shono.eu.auth0.com/oauth/token"
	payload := strings.NewReader(fmt.Sprintf("{\"client_id\":\"%s\",\"client_secret\":\"%s\",\"audience\":\"dev-api.shono.io\",\"grant_type\":\"client_credentials\"}", shonoKey, shonoSecret))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to authenticate with shono: %w", err)
	}

	defer res.Body.Close()
	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return "", err
	}

	t, fnd := body["access_token"]
	if !fnd {
		return "", fmt.Errorf("unable to find access token in response")
	}

	return t.(string), nil
}
