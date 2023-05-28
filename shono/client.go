package shono

import (
	"encoding/json"
	"fmt"
	"github.com/memphisdev/memphis.go"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/backbone"
	memphis2 "github.com/shono-io/go-shono/memphis"
	"io"
	"net/http"
	"strings"
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

func NewClient(id string, opts ...ClientOpt) (*Client, error) {
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
		return nil, fmt.Errorf("failed to get Config from shono: %w", err)
	}

	if cfg.Stream == nil {
		return nil, fmt.Errorf("no stream Config found")
	}

	c, err := memphis.Connect(cfg.Stream.Host, fmt.Sprintf("org.%s", co.Org), memphis.Password(cfg.Stream.Token))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the shono stream: %w", err)
	}

	return &Client{id: id, config: cfg, c: c}, nil
}

type Client struct {
	id     string
	config *Config
	c      *memphis.Conn

	run *memphis2.Runner
}

func (c *Client) Close() {
	if c.run != nil {
		c.run.Close()
	}

	c.c.Close()
}

func (c *Client) Listen(r *go_shono.Router) error {
	if c.run != nil {
		return fmt.Errorf("already listening")
	}

	c.run = memphis2.NewRunner(c.id, r, c.c)
	return c.run.Run()
}

func (c *Client) Id() string {
	return c.id
}

func (c *Client) Backbone() (backbone.Backbone, error) {
	if c.config.Backbone == nil {
		return nil, fmt.Errorf("no backbone Config found")
	}

	return backbone.NewBackbone(c.id, c.config.Backbone.Kind, c.config.Backbone.Properties)
}

func getConfig(at string, org string, baseUrl string) (*Config, error) {
	url := fmt.Sprintf("%s/organizations/%s/config", baseUrl, org)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", at))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to make the request to %s: %w", url, err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to %s did not return ok: %s", url, res.Status)
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	var body Config
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
