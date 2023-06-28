package runtime

import (
	"github.com/shono-io/shono/artifacts"
	"gopkg.in/yaml.v3"
)

func GenerateBenthosConfig(artifact *artifacts.Artifact, loglevel string) ([]byte, error) {
	inp := map[string]any{
		artifact.Input.Kind: artifact.Input.Config,
	}

	if artifact.Input.Logic != nil {
		inp["processors"] = artifact.Input.Logic.Steps
	}

	result := map[string]any{
		"input": inp,
		"pipeline": map[string]any{
			"processors": artifact.Logic.Steps,
		},
		"output": map[string]any{
			"switch": map[string]any{
				"cases": []map[string]any{
					{
						"check": "errored()",
						"output": map[string]any{
							artifact.DLQ.Kind: artifact.DLQ.Config,
						},
					},
					{
						"output": map[string]any{
							artifact.Output.Kind: artifact.Output.Config,
						},
					},
				},
			},
		},
		"logger": map[string]any{
			"level": loglevel,
		},
	}

	if artifact.Logic.Tests != nil {
		result["tests"] = artifact.Logic.Tests
	}

	return yaml.Marshal(result)
}
