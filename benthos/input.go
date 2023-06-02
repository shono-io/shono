package benthos

import (
	"github.com/shono-io/shono"
)

func (g *Generator) generateInput(result map[string]any, scope shono.Scope) (err error) {
	// -- get the list of events from all reaktors within the scope
	var events []shono.EventId
	for _, reaktor := range scope.Reaktors() {
		events = append(events, reaktor.InputEvent())
	}

	result["input"], err = g.bb.GetConsumerConfig(g.group, events)
	return err
}
