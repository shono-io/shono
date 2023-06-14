package benthos

import "github.com/shono-io/shono/core"

func generateSystemInput(env core.Environment, systemRef core.Reference, config map[string]any) (map[string]any, error) {
	system, err := env.ResolveSystem(systemRef)
	if err != nil {
		return nil, err
	}

	return system.AsInput(config)
}

func generateSystemOutput(env core.Environment, systemRef core.Reference, config map[string]any) (map[string]any, error) {
	system, err := env.ResolveSystem(systemRef)
	if err != nil {
		return nil, err
	}

	return system.AsOutput(config)
}
