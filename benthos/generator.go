package benthos

import (
	"context"
	"fmt"
	"github.com/shono-io/shono/backbone"
	"github.com/shono-io/shono/graph"
	"gopkg.in/yaml.v3"
	"os"
	"runtime/debug"
	"strings"
)

type GeneratorOpt func(g *Generator)

func WithThreads(threads int) GeneratorOpt {
	return func(g *Generator) {
		g.threads = threads
	}
}

func WithOutputDir(dir string) GeneratorOpt {
	return func(g *Generator) {
		g.outputDir = dir
	}
}

func WithGroup(group string) GeneratorOpt {
	return func(g *Generator) {
		g.group = group
	}
}

func NewGenerator(bb backbone.Backbone, opts ...GeneratorOpt) *Generator {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic("failed to read build info; maybe your project is not a go module?")
	}

	res := &Generator{
		group:     bi.Path,
		bb:        bb,
		threads:   1,
		outputDir: "generated_output",
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

type Generator struct {
	group     string
	bb        backbone.Backbone
	outputDir string
	threads   int
}

func (g *Generator) Generate(ctx context.Context, env graph.Environment) (err error) {
	// -- create the output directory
	if _, err := os.Stat(g.outputDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(g.outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}
		}
	}

	output := &GeneratorOutput{g.outputDir}

	// -- list through the scopes within the environment, generating each of them
	scopes, err := env.ListScopes()
	if err != nil {
		return fmt.Errorf("failed to list scopes: %w", err)
	}
	for _, scope := range scopes {
		if err := g.generateScope(ctx, output, env, scope); err != nil {
			return fmt.Errorf("failed to generate scope: %w", err)
		}
	}

	return nil
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

	return out.Write(fmt.Sprintf("%s-%s.yaml", scope.Key().Code(), concept.Key().Code()), result)
}

type GeneratorOutput struct {
	workDir string
}

func (g *GeneratorOutput) Write(path string, data map[string]any) (err error) {
	if !strings.HasSuffix(path, ".yml") && !strings.HasSuffix(path, ".yaml") {
		return fmt.Errorf("path must end with .yml or .yaml")
	}

	// -- generate the yaml based on the data
	b, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}

	// -- create the output file
	return os.WriteFile(fmt.Sprintf("%s/%s", g.workDir, path), b, 0644)
}
