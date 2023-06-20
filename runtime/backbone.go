package runtime

import "github.com/shono-io/shono/inventory"

type Backbone interface {
	EventLogName(event inventory.Event) string
	AsInput(id string, events ...inventory.Event) (map[string]any, error)
	AsOutput(events ...inventory.Event) (map[string]any, error)
}
