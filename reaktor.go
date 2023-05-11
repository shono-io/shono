package go_shono

import (
	"context"
	"errors"
	"time"
)

type ReaktorContext interface {
	Context

	Key() string
	Timestamp() time.Time
	Header(key string) []byte
}

type ReaktorFunc func(ctx context.Context, evt any, w Writer)

func MustNewReaktor(id string, opt ...ReaktorOpt) Reaktor {
	r, err := NewReaktor(id, opt...)
	if err != nil {
		panic(err)
	}

	return *r
}

func NewReaktor(id string, opt ...ReaktorOpt) (*Reaktor, error) {
	r := &Reaktor{
		Id: id,
	}

	for _, o := range opt {
		o(r)
	}

	if r.Handler == nil {
		return nil, errors.New("reaktor must have a handler")
	}

	if len(r.Listen) == 0 {
		return nil, errors.New("reaktor must listen for at least one event kind")
	}

	return r, nil
}

type ReaktorOpt func(r *Reaktor)

func ListenFor(events ...*EventMeta) ReaktorOpt {
	return func(r *Reaktor) {
		r.Listen = events
	}
}

func WithHandler(handler ReaktorFunc) ReaktorOpt {
	return func(r *Reaktor) {
		r.Handler = handler
	}
}

type Reaktor struct {
	Id      string
	Listen  []*EventMeta
	Handler ReaktorFunc
}
