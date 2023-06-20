package runtime

import (
	"github.com/shono-io/shono/artifacts"
	"gopkg.in/yaml.v3"
)

func GenerateBenthosConfig(artifact artifacts.Artifact) ([]byte, error) {

	result := map[string]any{
		"input": map[string]any{
			artifact.Input().Name: artifact.Input().Config,
		},
		"pipeline": map[string]any{
			"processors": artifact.Logic().Steps,
		},
		"output": map[string]any{
			"switch": map[string]any{
				"cases": []map[string]any{
					{
						"check": "errored()",
						"output": map[string]any{
							artifact.Error().Name: artifact.Error().Config,
						},
					},
					{
						"output": map[string]any{
							artifact.Output().Name: artifact.Output().Config,
						},
					},
				},
			},
		},
	}

	if artifact.Logic().Tests != nil {
		result["tests"] = artifact.Logic().Tests
	}

	return yaml.Marshal(result)
}
