package benthos

import (
	"context"
	"fmt"
	"github.com/shono-io/shono/graph"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"runtime/debug"
)

type GeneratorOpt func(g *Generator)

func WithThreads(threads int) GeneratorOpt {
	return func(g *Generator) {
		g.threads = threads
	}
}

func WithGroup(group string) GeneratorOpt {
	return func(g *Generator) {
		g.group = group
	}
}

func NewGenerator(opts ...GeneratorOpt) *Generator {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic("failed to read build info; maybe your project is not a go module?")
	}

	res := &Generator{
		group:   bi.Path,
		threads: 1,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

type Generator struct {
	group   string
	threads int
}

func (g *Generator) Generate(ctx context.Context, env graph.Environment) (output *GeneratorOutput, err error) {
	output = &GeneratorOutput{
		Streams: make([]Stream, 0),
	}

	// -- list through the scopes within the environment, generating each of them
	scopes, err := env.ListScopes()
	if err != nil {
		return nil, fmt.Errorf("failed to list scopes: %w", err)
	}
	for _, scope := range scopes {
		if err := g.generateScope(ctx, output, env, scope); err != nil {
			return nil, fmt.Errorf("failed to generate scope: %w", err)
		}
	}

	return output, nil
}

func (g *Generator) generateScope(ctx context.Context, out *GeneratorOutput, env graph.Environment, scope graph.Scope) (err error) {
	concepts, err := env.ListConceptsForScope(scope.Key())
	if err != nil {
		return fmt.Errorf("failed to list concepts for scope %q: %w", scope.Key().Code(), err)
	}

	for _, concept := range concepts {
		if err := g.generateConcept(ctx, out, env, scope, concept); err != nil {
			return fmt.Errorf("failed to generate concept %q: %w", concept.Key().Code(), err)
		}
	}

	return nil
}

func (g *Generator) generateConcept(ctx context.Context, out *GeneratorOutput, env graph.Environment, scope graph.Scope, concept graph.Concept) (err error) {
	reaktors, err := env.ListReaktorsForConcept(concept.Key())
	if err != nil {
		return fmt.Errorf("failed to list reaktors for concept: %w", err)
	}

	result := map[string]any{}
	if err := g.generateInput(ctx, result, env, scope, concept, reaktors); err != nil {
		return fmt.Errorf("failed to generate inputs: %w", err)
	}

	if err := g.generatePipeline(ctx, result, env, scope, concept, reaktors); err != nil {
		return fmt.Errorf("failed to generate pipeline: %w", err)
	}

	if err := g.generateOutput(ctx, result, env, scope, concept, reaktors); err != nil {
		return fmt.Errorf("failed to generate outputs: %w", err)
	}

	if err := g.generateCaches(ctx, result, env, scope, concept, reaktors); err != nil {
		return fmt.Errorf("failed to generate caches: %w", err)
	}

	if err := g.generateTests(ctx, result, env, scope, concept, reaktors); err != nil {
		return fmt.Errorf("failed to generate tests: %w", err)
	}

	out.RegisterUnit(concept, result)

	return nil
}

type Stream struct {
	Concept graph.Concept
	Unit    map[string]any
}

func (s *Stream) Dump(dir string) error {
	filename := fmt.Sprintf("%s.yaml", s.Concept.Key().CodeString())

	b, err := yaml.Marshal(s.Unit)
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/%s", dir, filename), b, 0644)
}

type GeneratorOutput struct {
	Streams []Stream
}

func (g *GeneratorOutput) RegisterUnit(concept graph.Concept, unit map[string]any) {
	g.Streams = append(g.Streams, Stream{
		Concept: concept,
		Unit:    unit,
	})
}

func (g *GeneratorOutput) Dump(dir string) error {
	// -- make sure the directory exists
	if _, err := os.Stat(dir); err == fs.ErrNotExist {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %q: %w", dir, err)
		}
	}

	for _, stream := range g.Streams {
		if err := stream.Dump(dir); err != nil {
			return fmt.Errorf("failed to dump stream: %w", err)
		}
	}

	return nil
}
