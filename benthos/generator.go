package benthos

import (
	"fmt"
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/shono-io/shono"
	"gopkg.in/yaml.v3"
	"io"
)

func NewGenerator(group string, bb shono.Backbone, threads int) *Generator {
	return &Generator{
		group:   group,
		bb:      bb,
		threads: threads,
	}
}

type Generator struct {
	group   string
	bb      shono.Backbone
	threads int
}

func (g *Generator) Generate(scope shono.Scope) (output shono.GeneratorOutput, err error) {
	result := map[string]any{}

	if err := g.generateInput(result, scope); err != nil {
		return nil, fmt.Errorf("failed to generate inputs: %w", err)
	}

	if err := g.generatePipeline(result, scope); err != nil {
		return nil, fmt.Errorf("failed to generate pipeline: %w", err)
	}

	if err := g.generateOutput(result, scope); err != nil {
		return nil, fmt.Errorf("failed to generate outputs: %w", err)
	}

	if err := g.generateCaches(result, scope); err != nil {
		return nil, fmt.Errorf("failed to generate caches: %w", err)
	}

	// -- convert the result to yaml
	yml, err := yaml.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal yaml: %w", err)
	}

	return &GeneratorOutput{yml}, nil
}

type GeneratorOutput struct {
	yml []byte
}

func (g *GeneratorOutput) Validate() error {
	b := service.NewStreamBuilder()
	if err := b.SetYAML(string(g.yml)); err != nil {
		return err
	}

	_, err := b.Build()
	return err
}

func (g *GeneratorOutput) Write(w io.Writer, opts ...shono.WriterOpt) (err error) {
	cfg := shono.NewWriterConfig(opts...)

	result := g.yml
	if cfg.Secure() {
		b := service.NewStreamBuilder()
		if err := b.SetYAML(string(g.yml)); err != nil {
			return err
		}
		res, err := b.AsYAML()
		if err != nil {
			return err
		}

		result = []byte(res)
	}

	_, err = w.Write(result)
	return err
}
