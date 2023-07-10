package runtime

import (
	"github.com/shono-io/shono/artifacts"
	"gopkg.in/yaml.v3"
	"strings"
)

func GenerateBenthosConfig(artifact *artifacts.Artifact, loglevel string, debug bool) ([]byte, error) {

	inp := map[string]any{}
	if debug {
		inp["stdin"] = map[string]any{
			"codec": "lines",
		}
		inp["processors"] = []map[string]any{
			{"mapping": strings.TrimSpace(`
meta = this.meta
root = this.root
`)},
		}
	} else {
		inp[artifact.Input.Kind] = artifact.Input.Config
		if artifact.Input.Logic != nil {
			inp["processors"] = artifact.Input.Logic.Steps
		}
	}

	out := map[string]any{}
	if debug {
		out["stdout"] = map[string]any{
			"codec": "lines",
		}
	} else {
		out[artifact.Output.Kind] = artifact.Output.Config
	}

	dlq := map[string]any{}
	if debug {
		dlq["stdout"] = map[string]any{
			"codec": "lines",
		}
	} else {
		dlq[artifact.DLQ.Kind] = artifact.DLQ.Config
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
						"check":  "errored()",
						"output": dlq,
					},
					{
						"output": out,
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
