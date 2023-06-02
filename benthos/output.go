package benthos

import "github.com/shono-io/shono"

func (g *Generator) generateOutput(result map[string]any, scope shono.Scope) (err error) {
	// -- get the list of events from the reaktors
	var events []shono.EventId
	for _, reaktor := range scope.Reaktors() {
		events = append(events, reaktor.OutputEvents()...)
	}

	result["output"], err = g.bb.GetProducerConfig(events)
	return err
}
