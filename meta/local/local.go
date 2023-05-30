package local

import (
	"github.com/shono-io/shono"
	"github.com/shono-io/shono/meta"
)

func NewLocalClient() meta.Client {
	return &localClient{
		scopeRepo: scopeRepo{
			scopes: make(map[string]shono.Scope),
		},
		conceptRepo: conceptRepo{
			concepts: make(map[string]shono.Concept),
		},
		reaktorRepo: reaktorRepo{
			reaktors: make(map[string]shono.Reaktor),
		},
		eventRepo: eventRepo{
			events: make(map[string]shono.Event),
		},
	}
}

type localClient struct {
	scopeRepo
	conceptRepo
	reaktorRepo
	eventRepo
}

func (l *localClient) Close() {
	// -- noop
}
