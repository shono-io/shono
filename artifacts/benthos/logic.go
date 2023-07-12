package benthos

import (
	"fmt"
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/inventory"
)

func generateLogic(logic inventory.Logic) (*artifacts.GeneratedLogic, error) {
	result := &artifacts.GeneratedLogic{
		Steps: make([]map[string]any, 0),
	}

	if len(logic.Steps()) == 0 {
		return nil, fmt.Errorf("logic has no steps")
	}

	for idx, step := range logic.Steps() {
		m, err := step.Build().MarshalBenthos(fmt.Sprintf("steps[%d]", idx))

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
	batch := map[string]any{}
	for _, ass := range test.Assertions {
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
		"name":        test.Summary,
		"environment": test.EnvironmentVars,
		"input_batch": []map[string]any{
			{
				"json_content": test.Input.Content,
				"metadata":     test.Input.Metadata,
			},
		},
		"output_batches": [][]map[string]any{
			{
				batch,
			},
		},
	}, nil
}
