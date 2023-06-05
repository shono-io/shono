package benthos

import (
	"context"
	"fmt"
	"github.com/shono-io/shono/graph"
	"strings"
)

func (g *Generator) generatePipeline(ctx context.Context, result map[string]any, env graph.Environment, scope graph.Scope, concept graph.Concept, reaktors []graph.Reaktor) (err error) {
	result["pipeline"] = map[string]any{
		"threads": g.threads,
	}

	return g.generateProcessors(ctx, result["pipeline"].(map[string]any), env, scope, concept, reaktors)
}

func (g *Generator) generateProcessors(ctx context.Context, result map[string]any, env graph.Environment, scope graph.Scope, concept graph.Concept, reaktors []graph.Reaktor) (err error) {
	var cases []map[string]any

	for _, reaktor := range reaktors {
		c, err := toCase(ctx, env, reaktor)
		if err != nil {
			return fmt.Errorf("failed to convert reaktor to case: %w", err)
		}

		if c != nil {
			cases = append(cases, c)
		}
	}

	result["processors"] = []map[string]any{
		{"switch": cases},
	}

	return nil
}

func toCase(ctx context.Context, env graph.Environment, reaktor graph.Reaktor) (map[string]any, error) {
	res := map[string]any{}

	var processors []map[string]any
	for _, e := range reaktor.Logic() {
		processor, err := toProcessor(ctx, env, e)
		if err != nil {
			return nil, fmt.Errorf("failed to generate processor: %w", err)
		}

		processors = append(processors, processor)
	}

	// -- we will add a check on the type header
	res["check"] = fmt.Sprintf("@io_shono_kind == %q", reaktor.InputEventKey().CodeString())
	res["processors"] = []map[string]any{
		{
			"label": labelize(reaktor.Key().CodeString()),
			"try":   processors,
		},
	}

	return res, nil
}

func toProcessor(ctx context.Context, env graph.Environment, l graph.Logic) (map[string]any, error) {
	switch lt := l.(type) {
	case graph.CatchLogic:
		var elements []map[string]any
		for _, el := range lt.Logics {
			element, err := toProcessor(ctx, env, el)
			if err != nil {
				return nil, err
			}

			elements = append(elements, element)
		}

		result := map[string]any{
			"catch": elements,
		}

		return result, nil

	case graph.ConditionalLogic:
		result := make(map[string]any)

		var cases []map[string]any
		for _, c := range lt.Cases {
			el, err := toProcessor(ctx, env, c)
			if err != nil {
				return nil, err
			}

			cases = append(cases, el)
		}

		if len(cases) > 0 {
			result["cases"] = cases
		}

		return map[string]any{
			"switch": cases,
		}, nil
	case graph.CaseLogic:
		var processors []map[string]any
		for _, b := range lt.Logics {
			// -- generate the element content
			el, err := toProcessor(ctx, env, b)
			if err != nil {
				return nil, err
			}

			processors = append(processors, el)
		}

		return map[string]any{
			"check":      lt.Check,
			"processors": processors,
		}, nil
	case graph.LogLogic:
		result := map[string]any{
			"level":   lt.Level,
			"message": lt.Message,
		}

		if lt.Mappings != nil {
			mappings, err := mappingsToString(ctx, lt.Mappings)
			if err != nil {
				return nil, err
			}

			result["fields_mapping"] = mappings
		}

		return map[string]any{
			"log": result,
		}, nil
	case graph.TransformLogic:
		mappings, err := mappingsToString(ctx, lt.Mappings)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"mapping": mappings,
		}, nil
	case graph.StoreLogic:
		// -- get the store
		store, err := env.GetStore(lt.StoreKey)
		if err != nil {
			return nil, err
		}

		if store == nil {
			return nil, fmt.Errorf("store %q not found", lt.StoreKey)
		}

		// -- get the storage linked to the store
		storage, err := env.GetStorage(store.StorageKey())
		if err != nil {
			return nil, err
		}

		if storage == nil {
			return nil, fmt.Errorf("storage %q not found", lt.StoreKey)
		}

		result := storage.Config()
		result["operation"] = lt.Operation
		result["collection"] = store.Collection()

		if lt.Key != nil {
			result["key"] = lt.Key
		}

		if len(lt.Value) > 0 {
			mappings, err := mappingsToString(ctx, lt.Value)
			if err != nil {
				return nil, err
			}

			result["value"] = mappings
		}

		if len(lt.Filters) > 0 {
			result["filters"] = lt.Filters
		}

		return map[string]any{
			storage.Kind(): result,
		}, nil
	default:
		return nil, fmt.Errorf("unknown logic type: %T", lt)
	}
}

func mappingsToString(ctx context.Context, mappings []graph.Mapping) (string, error) {
	var result []string
	for _, m := range mappings {
		mapping, err := m.Generate(ctx)
		if err != nil {
			return "", err
		}

		result = append(result, mapping)
	}

	return strings.Join(result, "\n"), nil
}
