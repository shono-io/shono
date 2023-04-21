package app

import (
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

type AppOpt func(a *App)

func WithSchemaRegistryOpts(opts ...sr.Opt) AppOpt {
	return func(a *App) {
		a.srOpts = append(a.srOpts, opts...)
	}
}

func WithKafkaOpts(opts ...kgo.Opt) AppOpt {
	return func(a *App) {
		a.kafkaOpts = append(a.kafkaOpts, opts...)
	}
}

func WithModules(modules ...Module) AppOpt {
	return func(a *App) {
		a.modules = append(a.modules, modules...)
	}
}
