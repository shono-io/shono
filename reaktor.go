package shono

import (
	"github.com/shono-io/shono/logic"
	"strings"
)

// == ENTITY ==========================================================================================================

type Reaktor interface {
	Entity
	InputEvent() EventId
	OutputEvents() []EventId
	Logic() logic.Logic
	Stores() []Store
	Tests() []ReaktorTest
}

type ReaktorTest interface {
	Summary() string
	Environment() map[string]any
	Mocks() map[string]any
	When() ReaktorTestEvent
	Then() []ReaktorTestCondition
}

type ReaktorTestEvent interface {
	Metadata() map[string]any
	Content() any
}

type ReaktorTestCondition interface {
	Kind() string
	Condition() any
}

func NewReaktor(scopeKey Key, code string, inputEvent EventId, logic logic.Logic, opts ...ReaktorOpt) Reaktor {
	result := &reaktor{
		key:          scopeKey.Child("reaktor", code),
		name:         strings.ToTitle(code),
		description:  "",
		inputEvent:   inputEvent,
		outputEvents: []EventId{},
		stores:       []Store{},
		tests:        []ReaktorTest{},
		logic:        logic,
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type reaktor struct {
	key          Key
	name         string
	description  string
	inputEvent   EventId
	outputEvents []EventId
	logic        logic.Logic
	stores       []Store
	tests        []ReaktorTest
}

func (r *reaktor) Key() Key {
	return r.key
}

func (r *reaktor) Name() string {
	return r.name
}

func (r *reaktor) Description() string {
	return r.description
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

func (r *reaktor) Tests() []ReaktorTest {
	return r.tests
}
