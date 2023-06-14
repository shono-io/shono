package runtime

import (
	"context"
	"fmt"
	"github.com/shono-io/shono/core"
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/graph"
	"strings"
)

func (g *Generator) generatePipeline(ctx context.Context, result map[string]any, reg graph.Registry, scope core.Scope, concept core.Concept, reaktors []graph.Reaktor) (err error) {
	result["pipeline"] = map[string]any{
		"threads": g.threads,
	}

	return g.generateProcessors(ctx, result["pipeline"].(map[string]any), reg, scope, concept, reaktors)
}

func (g *Generator) generateProcessors(ctx context.Context, result map[string]any, reg graph.Registry, scope core.Scope, concept core.Concept, reaktors []graph.Reaktor) (err error) {
	var cases []map[string]any

	for _, reaktor := range reaktors {
		c, err := toCase(ctx, reg, reaktor)
		if err != nil {
			return fmt.Errorf("failed to convert reaktor to case: %w", err)
		}

		if c != nil {
			cases = append(cases, c)
		}
	}

	cases = append(cases, map[string]any{
		"processors": []map[string]any{
			{
				"log": map[string]any{
					"level":   "TRACE",
					"message": "no processor for ${!this.format_json()}",
				},
			},
		},
	})

	result["processors"] = []map[string]any{
		{"switch": cases},
	}

	return nil
}

func toCase(ctx context.Context, reg graph.Registry, reaktor graph.Reaktor) (map[string]any, error) {
	res := map[string]any{}

	var processors []map[string]any
	for _, e := range reaktor.Logic {
		processor, err := toProcessor(ctx, reg, e)
		if err != nil {
			return nil, fmt.Errorf("failed to generate processor: %w", err)
		}

		processors = append(processors, processor)
	}

	// -- we will add a check on the type header
	res["check"] = fmt.Sprintf("@io_shono_kind == %q", reaktor.Input)
	res["processors"] = processors

	return res, nil
}

func toProcessor(ctx context.Context, reg graph.Registry, l graph.Logic) (map[string]any, error) {
	switch lt := l.(type) {
	case dsl.CatchLogicStep:
		var elements []map[string]any
		for _, el := range lt.Logics {
			element, err := toProcessor(ctx, reg, el)
			if err != nil {
				return nil, err
			}

			elements = append(elements, element)
		}

		result := map[string]any{
			"catch": elements,
		}

		return result, nil

	case dsl.ConditionalLogicStep:
		result := make(map[string]any)

		var cases []map[string]any
		for _, c := range lt.Cases {
			el, err := toProcessor(ctx, reg, c)
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
	case dsl.ConditionalCase:
		var processors []map[string]any
		for _, b := range lt.Logics {
			// -- generate the element content
			el, err := toProcessor(ctx, reg, b)
			if err != nil {
				return nil, err
			}

			processors = append(processors, el)
		}

		res := map[string]any{
			"processors": processors,
		}

		if lt.Check != "" {
			res["check"] = string(lt.Check)
		}

		return res, nil
	case dsl.LogLogicStep:
		result := map[string]any{
			"level":   string(lt.Level),
			"message": string(lt.Message),
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
	case dsl.TransformLogicStep:
		mappings, err := mappingsToString(ctx, lt.Mappings)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"mapping": mappings,
		}, nil
	case dsl.StoreLogicStep:
		result := map[string]any{
			"concept":   lt.Concept.String(),
			"operation": lt.Operation,
		}

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
			"store": result,
		}, nil
	default:
		return nil, fmt.Errorf("unknown logic type: %T", lt)
	}
}

func mappingsToString(ctx context.Context, mappings []dsl.Mapping) (string, error) {
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
