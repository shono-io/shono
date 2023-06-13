package graph

import (
	"context"
)

type LogStrategy string

var (
	PerScopeLogStrategy LogStrategy = "per_scope"
)

type Backbone interface {
	GetConsumerConfig(id string, events []EventReference) (map[string]any, error)
	GetProducerConfig(events []EventReference) (map[string]any, error)
	GetClient() (BackboneClient, error)
}

type BackboneClient interface {
	Produce(ctx context.Context, event EventReference, key string, payload map[string]any) error
}

type BackboneRepo interface {
	GetBackbone() (Backbone, error)
}
