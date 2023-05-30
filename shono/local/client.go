package local

import (
	"fmt"
	_ "github.com/benthosdev/benthos/v4/public/components/all"
	"github.com/shono-io/go-shono/shono"
)

type Config struct {
	Org    string `json:"org"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Url    string `json:"url"`
}

type ClientOpt func(*Config)

func NewClient(id string, scopeRepo shono.ScopeRepo, resourceRepo shono.ResourceRepo, opts ...ClientOpt) (shono.Client, error) {
	co := Config{
		Url: "https://dev-api.shono.io",
	}

	for _, opt := range opts {
		opt(&co)
	}

	cl := &client{id, co, scopeRepo, resourceRepo}

	if err := cl.validate(); err != nil {
		return nil, err
	}

	return cl, nil
}

type client struct {
	id  string
	cfg Config

	shono.ScopeRepo
	shono.ResourceRepo
}

func (c *client) validate() error {
	if c.id == "" {
		return fmt.Errorf("client id is not set")
	}

	if c.ScopeRepo == nil {
		return fmt.Errorf("client scope repo is not set")
	}

	if c.ResourceRepo == nil {
		return fmt.Errorf("client resource repo is not set")
	}

	return nil
}

func (c *client) Close() {}
