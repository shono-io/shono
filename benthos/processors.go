package benthos

import (
	"fmt"
	"github.com/shono-io/shono"
)

func (g *Generator) generateProcessors(result map[string]any, scope shono.Scope) (err error) {
	var cases []map[string]any

	for _, reaktor := range scope.Reaktors() {
		c, err := toCase(reaktor)
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

func toCase(reaktor shono.Reaktor) (map[string]any, error) {
	res := map[string]any{}

	processor, err := reaktor.Logic().Processor()
	if err != nil {
		return nil, fmt.Errorf("failed to get processor from logic: %w", err)
	}
	if len(processor) == 0 {
		return nil, nil
	}

	// -- we will add a check on the type header
	res["check"] = fmt.Sprintf("@io_shono_kind == %q", reaktor.InputEvent())
	res["processors"] = []map[string]any{
		{
			"label": labelize(reaktor.Key().CodeString()),
			"try": []map[string]any{
				processor,
			},
		},
	}

	return res, nil
}
