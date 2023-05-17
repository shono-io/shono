package events

import go_shono "github.com/shono-io/go-shono"

var (
	ConceptOperationFailed = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Concept", "OperationFailed"),
		new(ConceptOperationFailedEvent), go_shono.JsonSerde())
	ConceptCreationRequested = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Concept", "CreationRequested"),
		new(ConceptCreationRequestedEvent), go_shono.JsonSerde())
	ConceptCreated = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Concept", "Created"),
		new(ConceptCreatedEvent), go_shono.JsonSerde())
	ConceptDeletionRequested = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Concept", "DeletionRequested"),
		new(ConceptDeletionRequestedEvent), go_shono.JsonSerde())
	ConceptDeleted = go_shono.NewEvent(
		go_shono.NewEventId("io.shono", "core", "Concept", "Deleted"),
		new(ConceptDeletedEvent), go_shono.JsonSerde())
)

type ConceptOperationFailedEvent struct {
	Reason string `json:"reason"`
}

type ConceptCreationRequestedEvent struct {
	Organization string `json:"organization"`
	Scope        string `json:"scope"`
	Code         string `json:"code"`
	Name         string `json:"name"`
}

type ConceptCreatedEvent struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type ConceptDeletionRequestedEvent struct {
	Organization string `json:"organization"`
	Scope        string `json:"scope"`
	Code         string `json:"code"`
}

type ConceptDeletedEvent struct {
	Organization string `json:"organization"`
	Scope        string `json:"scope"`
	Code         string `json:"code"`
}
