package events

import go_shono "github.com/shono-io/go-shono"

var (
	ReaktorOperationFailed = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Reaktor", "OperationFailed"),
		new(ReaktorOperationFailedEvent), go_shono.JsonSerde())
	ReaktorCreationRequested = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Reaktor", "CreationRequested"),
		new(ReaktorCreationRequestedEvent), go_shono.JsonSerde())
	ReaktorCreated = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Reaktor", "Created"),
		new(ReaktorCreatedEvent), go_shono.JsonSerde())
	ReaktorDeletionRequested = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Reaktor", "DeletionRequested"),
		new(ReaktorDeletionRequestedEvent), go_shono.JsonSerde())
	ReaktorDeleted = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Reaktor", "Deleted"),
		new(ReaktorDeletedEvent), go_shono.JsonSerde())
)

type ReaktorOperationFailedEvent struct {
	Reason string `json:"reason"`
}

type ReaktorCreationRequestedEvent struct {
	Organization string `json:"organization"`
	Scope        string `json:"scope"`
	Code         string `json:"code"`
	Name         string `json:"name"`
}

type ReaktorCreatedEvent struct {
}

type ReaktorDeletionRequestedEvent struct {
}

type ReaktorDeletedEvent struct {
}
