package benthos

import (
	"github.com/shono-io/shono/core"
)

func generateBackboneInput(env core.Environment, eventRefs []core.Reference) (map[string]any, error) {
	// -- get the backbone
	bb := env.Backbone()

	// -- resolve the events
	var err error
	events := make([]core.Event, len(eventRefs))
	for i, ref := range eventRefs {
		events[i], err = env.ResolveEvent(ref)
		if err != nil {
			return nil, err
		}
	}

	return bb.AsInput(env.ApplicationId(), events...)
}

func generateBackboneOutput(env core.Environment, eventRefs []core.Reference) (map[string]any, error) {
	// -- get the backbone
	bb := env.Backbone()

	// -- resolve the events
	var err error
	events := make([]core.Event, len(eventRefs))
	for i, ref := range eventRefs {
		events[i], err = env.ResolveEvent(ref)
		if err != nil {
			return nil, err
		}
	}

	return bb.AsOutput(events...)
}
