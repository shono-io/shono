package reaktors

import "github.com/shono-io/go-shono/events"

type ReaktorFunc func(ctx ReaktorContext, key string, event any)

type ReaktorOpt func(*ReaktorInfo)

func NewReaktorInfo(consumes events.Kind, fn ReaktorFunc, opts ...ReaktorOpt) *ReaktorInfo {
	reaktor := &ReaktorInfo{
		Consumes: consumes,
		Fn:       fn,
	}

	for _, opt := range opts {
		opt(reaktor)
	}

	return reaktor
}

type ReaktorInfo struct {
	Consumes events.Kind
	Fn       ReaktorFunc
}
