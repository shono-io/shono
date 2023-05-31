package shono

import (
	"github.com/shono-io/shono/logic"
)

// == ENTITY ==========================================================================================================

type Reaktor interface {
	Entity
	InputEvent() EventId
	OutputEvents() []EventId
	Logic() logic.Logic
	Stores() []Store
}

func NewReaktor(scopeKey Key, code string, inputEvent EventId, logic logic.Logic, opts ...ReaktorOpt) Reaktor {
	result := &reaktor{
		entity:       newEntity(scopeKey.Child("reaktor", code)),
		inputEvent:   inputEvent,
		outputEvents: []EventId{},
		stores:       []Store{},
		logic:        logic,
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type reaktor struct {
	*entity
	inputEvent   EventId
	outputEvents []EventId
	logic        logic.Logic
	stores       []Store
}

func (r *reaktor) InputEvent() EventId {
	return r.inputEvent
}

func (r *reaktor) OutputEvents() []EventId {
	return r.outputEvents
}

func (r *reaktor) Logic() logic.Logic {
	return r.logic
}

func (r *reaktor) Stores() []Store {
	return r.stores
}
