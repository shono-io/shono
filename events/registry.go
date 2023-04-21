package events

import "context"

type EventRegistry interface {
	Event(ctx context.Context, kind Kind) (*EventInfo, error)
}
