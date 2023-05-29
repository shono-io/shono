package shono

import (
	"context"
	"fmt"
)

// == ENTITY ==========================================================================================================

type EventId string

type Event interface {
	Entity
	GetSchema() string
	Id() EventId
}

func NewEvent(scopeCode, conceptCode, code, name, description, schema string) Event {
	return &event{
		NewEntity(fmt.Sprintf("%s:%s:%s", scopeCode, conceptCode, code), code, name, description),
		scopeCode,
		conceptCode,
		schema,
	}
}

type event struct {
	Entity
	scopeCode   string
	conceptCode string
	schema      string
}

func (e *event) GetSchema() string {
	if e == nil {
		return ""
	}

	return e.schema
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
