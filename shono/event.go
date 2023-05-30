package shono

import (
	"context"
	"fmt"
	"strings"
)

// == ENTITY ==========================================================================================================

type EventId string

func (e EventId) part(i int) string {
	parts := strings.Split(string(e), ":")
	if len(parts) != 3 {
		return ""
	}
	return parts[i]
}

func (e EventId) Scope() string {
	return e.part(0)
}

func (e EventId) Concept() string {
	return e.part(1)
}

func (e EventId) Code() string {
	return e.part(2)
}

type Event interface {
	Entity
	Id() EventId
}

func NewEvent(scopeCode, conceptCode, code, name, description string) Event {
	return &event{
		NewEntity(fmt.Sprintf("%s:%s:%s", scopeCode, conceptCode, code), code, name, description),
		scopeCode,
		conceptCode,
	}
}

type event struct {
	Entity
	scopeCode   string
	conceptCode string
}

func (e *event) Id() EventId {
	if e == nil {
		return ""
	}

	return EventId(e.FQN())
}

// == REPO ============================================================================================================

type EventRepo interface {
	GetEvent(ctx context.Context, code string) (Event, bool, error)
	AddEvent(ctx context.Context, event Event) error
	RemoveEvent(ctx context.Context, code string) error
}
