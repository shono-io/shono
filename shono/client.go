package shono

import "github.com/shono-io/go-shono/backbone"

type Client interface {
	ScopeRepo

	Backbone() backbone.Backbone

	Close()
}
