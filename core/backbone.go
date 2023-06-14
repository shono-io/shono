package core

type Backbone interface {
	EventLogName(event Event) string
	AsInput(id string, events ...Event) (map[string]any, error)
	AsOutput(events ...Event) (map[string]any, error)
}
