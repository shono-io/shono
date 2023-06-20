package benthos

import (
	"fmt"
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/inventory"
)

func generateLogic(logic inventory.Logic) (*artifacts.GeneratedLogic, error) {
	result := &artifacts.GeneratedLogic{
		Steps: make([]map[string]any, 0),
	}

	if len(logic.Steps()) == 0 {
		return nil, fmt.Errorf("logic has no steps")
	}

	for _, step := range logic.Steps() {
		m, err := generateLogicStep(step)

		if err != nil {
			return nil, err
		}

		result.Steps = append(result.Steps, m)
	}

	if logic.Tests() != nil {
		result.Tests = make([]map[string]any, 0)
		for _, test := range logic.Tests() {
			t, err := generateTest(test)
			if err != nil {
				return nil, err
			}

			result.Tests = append(result.Tests, t)
		}
	}

	return result, nil
}

func generateTest(test inventory.Test) (map[string]any, error) {
	var outputBatch []map[string]any
	for _, ass := range test.Assertions() {
		batch := map[string]any{}
		if ass.Payload() != nil {
			if ass.Strict() {
				batch["json_equals"] = ass.Payload()
			} else {
				batch["json_contains"] = ass.Payload()
			}
		} else if ass.Metadata() != nil {
			batch["metadata_equals"] = ass.Metadata()
		}
	}

	return map[string]any{
		"name":        test.Summary(),
		"environment": test.EnvironmentVars(),
		"input_batch": []map[string]any{
			{
				"content":  test.Input().Content,
				"metadata": test.Input().Metadata,
			},
		},
		"output_batches": [][]map[string]any{
			outputBatch,
		},
	}, nil
}

func generateLogicStep(step inventory.LogicStep) (map[string]any, error) {
	switch t := step.(type) {
	case dsl.CatchLogicStep:
		return generateCatchLogicStep(t)
	case dsl.LogLogicStep:
		return generateLogLogicStep(t)
	case dsl.ConditionalLogicStep:
		return generateConditionalLogicStep(t)
	case dsl.RawLogicStep:
		return generateRawLogicStep(t)
	case dsl.StoreLogicStep:
		return generateStoreLogicStep(t)
	case dsl.TransformLogicStep:
		return generateTransformLogicStep(t)
	default:
		return nil, fmt.Errorf("unknown logic step type: %T", t)
	}
}

func generateCatchLogicStep(step dsl.CatchLogicStep) (map[string]any, error) {
	var elements []map[string]any
	for _, el := range step.Steps {
		element, err := generateLogicStep(el)
		if err != nil {
			return nil, err
		}

		elements = append(elements, element)
	}

	result := map[string]any{
		"catch": elements,
	}

	return result, nil
}

func generateLogLogicStep(step dsl.LogLogicStep) (map[string]any, error) {
	result := map[string]any{
		"level":   string(step.Level),
		"message": string(step.Message),
	}

	if step.Mapping != nil {
		switch step.Mapping.Language {
		case "bloblang":
			result["fields_mapping"] = step.Mapping.Sourcecode
		default:
			return nil, fmt.Errorf("unknown mapping language: %s", step.Mapping.Language)
		}
	}

	return map[string]any{
		"log": result,
	}, nil
}

func generateConditionalLogicStep(step dsl.ConditionalLogicStep) (map[string]any, error) {
	result := make(map[string]any)

	var cases []map[string]any
	for _, c := range step.Cases {
		var processors []map[string]any
		for _, b := range c.Steps {
			el, err := generateLogicStep(b)
			if err != nil {
				return nil, err
			}

			processors = append(processors, el)
		}

		res := map[string]any{
			"processors": processors,
		}

		if c.Check != "" {
			res["check"] = c.Check
		}

		cases = append(cases, res)
	}

	if len(cases) > 0 {
		result["cases"] = cases
	}

	return map[string]any{
		"switch": cases,
	}, nil
}

func generateRawLogicStep(step dsl.RawLogicStep) (map[string]any, error) {
	return step.Content, nil
}

func generateStoreLogicStep(step dsl.StoreLogicStep) (map[string]any, error) {
	result := map[string]any{
		"concept":   step.Concept.String(),
		"operation": step.Operation,
	}

	if step.Key != "" {
		result["key"] = step.Key
	}

	if step.Value != nil {
		switch step.Value.Language {
		case "bloblang":
			result["value"] = step.Value.Sourcecode
		default:
			return nil, fmt.Errorf("unknown mapping language: %s", step.Value.Language)
		}
	}

	if len(step.Filters) > 0 {
		result["filters"] = step.Filters
	}

	return map[string]any{
		"store": result,
	}, nil
}

func generateTransformLogicStep(step dsl.TransformLogicStep) (map[string]any, error) {
	result := map[string]any{}

	switch step.Mapping.Language {
	case "bloblang":
		result["mapping"] = step.Mapping.Sourcecode
	default:
		return nil, fmt.Errorf("unknown mapping language: %s", step.Mapping.Language)
	}
	return result, nil
}
