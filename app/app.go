package app

import (
	"fmt"
	"github.com/shono-io/go-shono/handlers"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

func NewApp(domain string, id string, opt ...AppOpt) (*App, error) {
	result := &App{
		Identity: Identity{
			Domain: domain,
			Id:     id,
		},
		kafkaOpts: []kgo.Opt{
			kgo.ConsumerGroup(id),
		},
	}

	for _, o := range opt {
		o(result)
	}

	if err := result.validate(); err != nil {
		return nil, err
	}

	return result, nil
}

type Identity struct {
	Domain string
	Id     string
}

type App struct {
	Identity

	srOpts    []sr.Opt
	kafkaOpts []kgo.Opt

	modules []Module
}

func (a App) validate() error {
	// -- make sure we have a valid config
	if a.Identity.Domain == "" {
		return ErrMissingDomain
	}

	if a.Identity.Id == "" {
		return ErrMissingAppId
	}

	return nil
}

func (a App) Run() error {
	src, err := sr.NewClient(a.srOpts...)
	if err != nil {
		return fmt.Errorf("failed to create schema registry client: %w", err)
	}

	cats, err := NewCatalogs(src)
	cats.Handlers, err = handlers.NewCatalog()
	if err != nil {
		return fmt.Errorf("failed to create the catalogs: %w", err)
	}

	bb, err := NewBackbone(cats, a.kafkaOpts...)
	if err != nil {
		return fmt.Errorf("failed to create backbone: %w", err)
	}

	bb.Execute()

	return nil
}
