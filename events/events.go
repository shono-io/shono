package events

import go_shono "github.com/shono-io/go-shono"

var (
	EventOperationFailed = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Event", "OperationFailed"),
		new(EventOperationFailedEvent), go_shono.JsonSerde())
	EventCreationRequested = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Event", "CreationRequested"),
		new(EventCreationRequestedEvent), go_shono.JsonSerde())
	EventCreated = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Event", "Created"),
		new(EventCreatedEvent), go_shono.JsonSerde())
	EventDeletionRequested = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Event", "DeletionRequested"),
		new(EventDeletionRequestedEvent), go_shono.JsonSerde())
	EventDeleted = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Event", "Deleted"),
		new(EventDeletedEvent), go_shono.JsonSerde())
)

type EventOperationFailedEvent struct {
	Reason string `json:"reason"`
}

type EventCreationRequestedEvent struct {
	Organization string `json:"organization"`
	Scope        string `json:"scope"`
	Concept      string `json:"concept"`
	Code         string `json:"code"`
	Name         string `json:"name"`
}

type EventCreatedEvent struct {
	Key  string `json:"key"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type EventDeletionRequestedEvent struct {
	Organization string `json:"organization"`
	Scope        string `json:"scope"`
	Concept      string `json:"concept"`
	Code         string `json:"code"`
}

type EventDeletedEvent struct {
	Organization string `json:"organization"`
	Scope        string `json:"scope"`
	Concept      string `json:"concept"`
	Code         string `json:"code"`
}
