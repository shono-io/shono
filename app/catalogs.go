package app

import (
	"fmt"
	"github.com/shono-io/go-shono/events"
	"github.com/shono-io/go-shono/handlers"
	"github.com/shono-io/go-shono/states"
	"github.com/twmb/franz-go/pkg/sr"
)

func NewCatalogs(src *sr.Client) (*Catalogs, error) {
	evt, err := events.NewCatalog(src)
	if err != nil {
		return nil, fmt.Errorf("unable to create events catalog: %w", err)
	}

	han, err := handlers.NewCatalog()
	if err != nil {
		return nil, fmt.Errorf("unable to create handlers catalog: %w", err)
	}

	sta, err := states.NewCatalog()
	if err != nil {
		return nil, fmt.Errorf("unable to create states catalog: %w", err)
	}

	return &Catalogs{
		Events:   evt,
		Handlers: han,
		States:   sta,
	}, nil
}

type Catalogs struct {
	Events   *events.Catalog
	Handlers *handlers.Catalog
	States   *states.Catalog
}
