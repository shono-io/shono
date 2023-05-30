package shono

import (
	"context"
	"fmt"
	"github.com/shono-io/go-shono/shono/logic"
)

// == ENTITY ==========================================================================================================

type Reaktor interface {
	Entity
	InputEvent() EventId
	OutputEvents() []EventId
	Logic() logic.Logic
}

func NewReaktor(scopeCode, code, name, description string, inputEvent EventId, logic logic.Logic, outputEvents ...EventId) Reaktor {
	return &reaktor{
		ScopeCode:    scopeCode,
		Entity:       NewEntity(fmt.Sprintf("%s:%s", scopeCode, code), code, name, description),
		inputEvent:   inputEvent,
		outputEvents: outputEvents,
		logic:        logic,
	}
}

type reaktor struct {
	ScopeCode string
	Entity
	inputEvent   EventId
	outputEvents []EventId
	logic        logic.Logic
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

// == REPO ============================================================================================================

type ReaktorRepo interface {
	GetReaktor(ctx context.Context, code string) (Reaktor, bool, error)
	AddReaktor(ctx context.Context, reaktor Reaktor) error
	RemoveReaktor(ctx context.Context, code string) error
}
