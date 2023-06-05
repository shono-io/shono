package benthos

import (
	"context"
	"fmt"
	"github.com/shono-io/shono/graph"
)

func (g *Generator) generateTests(ctx context.Context, result map[string]any, env graph.Environment, scope graph.Scope, concept graph.Concept, reaktors []graph.Reaktor) (err error) {
	var tests []map[string]any

	for _, reaktor := range reaktors {
		rt, err := generateTestsForReaktor(ctx, env, reaktor)
		if err != nil {
			return fmt.Errorf("failed to convert reaktor to test: %w", err)
		}

		tests = append(tests, rt...)
	}

	result["tests"] = tests
	return nil
}

func generateTestsForReaktor(ctx context.Context, env graph.Environment, reaktor graph.Reaktor) (result []map[string]any, err error) {
	for _, test := range reaktor.Tests() {
		t, err := generateTest(ctx, test)
		if err != nil {
			return nil, fmt.Errorf("failed to generate test: %w", err)
		}

		result = append(result, t)
	}

	return result, nil
}

func generateTest(ctx context.Context, test graph.ReaktorTest) (result map[string]any, err error) {
	if test.Event == nil {
		return nil, fmt.Errorf("test event is nil")
	}

	var conditions []map[string]any
	for _, condition := range test.Conditions {
		condition, err := generateCondition(ctx, condition)
		if err != nil {
			return nil, fmt.Errorf("failed to generate condition: %w", err)
		}

		conditions = append(conditions, condition)
	}

	return map[string]any{
		"name":        test.Summary,
		"environment": test.Environment,
		"mocks":       test.Mocks,
		"input_batch": []map[string]any{
			{
				"metadata":     test.Event.Metadata,
				"json_content": test.Event.Content,
			},
		},
		"output_batches": [][]map[string]any{
			conditions,
		},
	}, nil
}

func generateCondition(ctx context.Context, condition graph.ReaktorTestCondition) (map[string]any, error) {
	switch c := condition.(type) {
	case graph.BloblangReaktorTestCondition:
		return map[string]any{
			"bloblang": c.Expression,
		}, nil
	case graph.MetadataReaktorTestCondition:
		kind := "metadata_contains"
		if c.Strict {
			kind = "metadata_equals"
		}

		return map[string]any{
			kind: c.Values,
		}, nil
	case graph.PayloadReaktorTestCondition:
		kind := "json_contains"
		if c.Strict {
			kind = "json_equals"
		}

		return map[string]any{
			kind: c.Values,
		}, nil
	default:
		return nil, fmt.Errorf("unknown condition type: %T", c)
	}
}
