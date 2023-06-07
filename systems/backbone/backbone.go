package backbone

import (
	"context"
	"github.com/shono-io/shono/commons"
)

type LogStrategy string

var (
	PerScopeLogStrategy LogStrategy = "per_scope"
)

type Backbone interface {
	GetConsumerConfig(id string, events []commons.Key) (map[string]any, error)
	GetProducerConfig(events []commons.Key) (map[string]any, error)
	GetClient() (Client, error)
}

type Client interface {
	Produce(ctx context.Context, event commons.Key, key commons.Key, payload map[string]any) error
}
