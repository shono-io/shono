package local

import (
	"fmt"
	"github.com/shono-io/go-shono/backbone"
	"github.com/shono-io/go-shono/shono"
)

type Config struct {
	Org    string `json:"org"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Url    string `json:"url"`
}

type ClientOpt func(*Config)

func NewClient(id string, bb backbone.Backbone, scopeRepo shono.ScopeRepo, opts ...ClientOpt) (shono.Client, error) {
	co := Config{
		Url: "https://dev-api.shono.io",
	}

	for _, opt := range opts {
		opt(&co)
	}

	cl := &client{id, co, bb, scopeRepo}

	if err := cl.validate(); err != nil {
		return nil, err
	}

	return cl, nil
}

type client struct {
	id  string
	cfg Config

	bb backbone.Backbone

	shono.ScopeRepo
}

func (c *client) validate() error {
	if c.id == "" {
		return fmt.Errorf("client id is not set")
	}

	if c.bb == nil {
		return fmt.Errorf("client backbone is not set")
	}

	if c.ScopeRepo == nil {
		return fmt.Errorf("client scope repo is not set")
	}

	return nil
}

func (c *client) Backbone() backbone.Backbone {
	return c.bb
}

func (c *client) Close() {}
