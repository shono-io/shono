package shono

import go_shono "github.com/shono-io/go-shono"

var (
	ScopeOperationFailed = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Scope", "OperationFailed"),
		new(ScopeOperationFailedEvent), go_shono.JsonSerde())
	ScopeCreationRequested = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Scope", "CreationRequested"),
		new(ScopeCreationRequestedEvent), go_shono.JsonSerde())
	ScopeCreated = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Scope", "Created"),
		new(ScopeCreatedEvent), go_shono.JsonSerde())
	ScopeDeletionRequested = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Scope", "DeletionRequested"),
		new(ScopeDeletionRequestedEvent), go_shono.JsonSerde())
	ScopeDeleted = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Scope", "Deleted"),
		new(ScopeDeletedEvent), go_shono.JsonSerde())
)

type ScopeOperationFailedEvent struct {
	Reason string `json:"reason"`
}

type ScopeCreationRequestedEvent struct {
	Organization string `json:"organization"`
	Code         string `json:"code"`
	Name         string `json:"name"`
}

type ScopeCreatedEvent struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ScopeDeletionRequestedEvent struct {
	Organization string `json:"organization"`
	Code         string `json:"code"`
}

type ScopeDeletedEvent struct {
	Organization string `json:"organization"`
	Code         string `json:"code"`
}
