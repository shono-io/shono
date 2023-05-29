package shono

import (
	"context"
	"fmt"
	"github.com/benthosdev/benthos/v4/public/service"
)

type BenthosConfigSection interface {
	YAML() string
}

type Input interface {
	BenthosConfigSection
}

type Output interface {
	BenthosConfigSection
}

type Processor interface {
	BenthosConfigSection
}

type BenthosOpt func(b *benthosReaktor)

func WithInput(input Input) BenthosOpt {
	return func(b *benthosReaktor) {
		b.inputs = append(b.inputs, input)
	}
}

func WithProcessor(processor Processor) BenthosOpt {
	return func(b *benthosReaktor) {
		b.processors = append(b.processors, processor)
	}
}

func WithOutput(output Output) BenthosOpt {
	return func(b *benthosReaktor) {
		b.outputs = append(b.outputs, output)
	}
}

func WithLogger(logger service.PrintLogger) BenthosOpt {
	return func(b *benthosReaktor) {
		b.logger = logger
	}
}

func WithThreads(threads int) BenthosOpt {
	return func(b *benthosReaktor) {
		b.threads = threads
	}
}

func NewBenthosReaktor(scopeCode, code, name, description string, opts ...BenthosOpt) Reaktor {
	result := &benthosReaktor{
		Entity:     NewEntity(fmt.Sprintf("%s:%s", scopeCode, code), code, name, description),
		inputs:     []Input{},
		processors: []Processor{},
		outputs:    []Output{},
		threads:    0,
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type benthosReaktor struct {
	Entity

	inputs     []Input
	outputs    []Output
	processors []Processor

	logger  service.PrintLogger
	threads int

	stream *service.Stream
}

func (b *benthosReaktor) Run(ctx context.Context) (err error) {
	// -- construct the stream
	builder := service.NewStreamBuilder()

	// -- add the inputs
	for _, input := range b.inputs {
		if err := builder.AddInputYAML(input.YAML()); err != nil {
			return err
		}
	}

	// -- add the processors
	for _, processor := range b.processors {
		if err := builder.AddProcessorYAML(processor.YAML()); err != nil {
			return err
		}
	}

	// -- add the outputs
	for _, output := range b.outputs {
		if err := builder.AddOutputYAML(output.YAML()); err != nil {
			return err
		}
	}

	// -- set the logger
	if b.logger != nil {
		builder.SetPrintLogger(b.logger)
	}

	// -- set the number of threads
	if b.threads > 0 {
		builder.SetThreads(b.threads)
	}

	b.stream, err = builder.Build()
	if err != nil {
		return err
	}

	// -- run the stream
	return b.stream.Run(ctx)
}

func (b *benthosReaktor) Close() error {
	if b.stream != nil {
		return b.stream.Stop(context.Background())
	}

	return nil
}
