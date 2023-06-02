package benthos

import "github.com/shono-io/shono"

func (g *Generator) generatePipeline(result map[string]any, scope shono.Scope) (err error) {
	result["pipeline"] = map[string]any{
		"threads": g.threads,
	}

	return g.generateProcessors(result["pipeline"].(map[string]any), scope)
}
