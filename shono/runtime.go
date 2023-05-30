package shono

import (
	"context"
	"errors"
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

type Opt func(*Runtime)

func WithId(id string) Opt {
	return func(r *Runtime) {
		r.id = id
	}
}

func WithBackbone(bb Backbone) Opt {
	return func(r *Runtime) {
		r.bb = bb
	}
}

func WithReaktor(reaktors ...Reaktor) Opt {
	return func(r *Runtime) {
		r.reaktors = append(r.reaktors, reaktors...)
	}
}

func WithLogger(logger service.PrintLogger) Opt {
	return func(r *Runtime) {
		r.logger = logger
	}
}

func WithThreads(threads int) Opt {
	return func(r *Runtime) {
		r.threads = threads
	}
}

func NewRuntime(opts ...Opt) (*Runtime, error) {
	result := &Runtime{
		id:       xid.New().String(),
		bb:       nil,
		reaktors: []Reaktor{},
		logger:   nil,
		threads:  1,
	}

	for _, opt := range opts {
		opt(result)
	}

	if err := result.validate(); err != nil {
		return nil, err
	}

	return result, nil
}

type Runtime struct {
	id       string
	bb       Backbone
	reaktors []Reaktor

	logger  service.PrintLogger
	threads int

	stream *service.Stream
}

var ErrMissingBackbone = errors.New("no backbone provided")
var ErrMissingReaktor = errors.New("no reaktor provided")

func (r *Runtime) validate() error {
	if r.bb == nil {
		return ErrMissingBackbone
	}

	if len(r.reaktors) == 0 {
		return ErrMissingReaktor
	}

	return nil
}

func (r *Runtime) Run(ctx context.Context) (err error) {
	builder, err := NewBuilder(r.id, r.bb, r.reaktors)
	if err != nil {
		return err
	}

	// -- set the logger
	if r.logger != nil {
		builder.SetPrintLogger(r.logger)
	}

	// -- set the number of threads
	if r.threads > 0 {
		builder.SetThreads(r.threads)
	}

	yml, _ := builder.AsYAML()
	logrus.Debugf("Benthos stream:\n%s", yml)

	r.stream, err = builder.Build()
	if err != nil {
		return err
	}

	// -- run the stream
	return r.stream.Run(ctx)
}

func (r *Runtime) Close() error {
	if r.stream != nil {
		return r.stream.Stop(context.Background())
	}

	return nil
}
