package benthos

import (
	"fmt"
	"github.com/shono-io/shono/core"
	"github.com/shono-io/shono/dsl"
)

func generateLogic(env core.Environment, logic core.Logic) ([]map[string]any, error) {
	var elements []map[string]any

	if len(logic.Steps()) == 0 {
		return nil, fmt.Errorf("logic has no steps")
	}

	for _, step := range logic.Steps() {
		m, err := generateLogicStep(env, step)

		if err != nil {
			return nil, err
		}

		elements = append(elements, m)
	}

	return elements, nil
}

func generateLogicStep(env core.Environment, step core.LogicStep) (map[string]any, error) {
	switch t := step.(type) {
	case dsl.CatchLogicStep:
		return generateCatchLogicStep(env, t)
	case dsl.LogLogicStep:
		return generateLogLogicStep(env, t)
	case dsl.ConditionalLogicStep:
		return generateConditionalLogicStep(env, t)
	case dsl.RawLogicStep:
		return generateRawLogicStep(env, t)
	case dsl.StoreLogicStep:
		return generateStoreLogicStep(env, t)
	case dsl.TransformLogicStep:
		return generateTransformLogicStep(env, t)
	default:
		return nil, fmt.Errorf("unknown logic step type: %T", t)
	}
}

func generateCatchLogicStep(env core.Environment, step dsl.CatchLogicStep) (map[string]any, error) {
	var elements []map[string]any
	for _, el := range step.Steps {
		element, err := generateLogicStep(env, el)
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

func generateLogLogicStep(env core.Environment, step dsl.LogLogicStep) (map[string]any, error) {
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

func generateConditionalLogicStep(env core.Environment, step dsl.ConditionalLogicStep) (map[string]any, error) {
	result := make(map[string]any)

	var cases []map[string]any
	for _, c := range step.Cases {
		var processors []map[string]any
		for _, b := range c.Steps {
			el, err := generateLogicStep(env, b)
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

func generateRawLogicStep(env core.Environment, step dsl.RawLogicStep) (map[string]any, error) {
	return step.Content, nil
}

func generateStoreLogicStep(env core.Environment, step dsl.StoreLogicStep) (map[string]any, error) {
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

func generateTransformLogicStep(env core.Environment, step dsl.TransformLogicStep) (map[string]any, error) {
	result := map[string]any{}

	switch step.Mapping.Language {
	case "bloblang":
		result["mapping"] = step.Mapping.Sourcecode
	default:
		return nil, fmt.Errorf("unknown mapping language: %s", step.Mapping.Language)
	}
	return result, nil
}
