package shono

import "io"

type Generator interface {
	Generate(scope Scope) (GeneratorOutput, error)
}

type GeneratorOutput interface {
	Validate() error
	Write(w io.Writer, opts ...WriterOpt) (err error)
}

func NewWriterConfig(opts ...WriterOpt) *WriterConfig {
	c := &WriterConfig{
		secure: true,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type WriterConfig struct {
	secure bool
}

func (c *WriterConfig) Secure() bool {
	return c.secure
}

type WriterOpt func(*WriterConfig)

func NonSecure() WriterOpt {
	return func(c *WriterConfig) {
		c.secure = false
	}
}
