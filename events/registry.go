package events

import (
	"context"
	"github.com/twmb/franz-go/pkg/kgo"
)

type EventRegistry interface {
	Event(ctx context.Context, kind Kind) (*EventInfo, error)

	MustAsRecord(ctx context.Context, kind Kind, key string, value any) *kgo.Record
	AsRecord(ctx context.Context, kind Kind, key string, value any) (*kgo.Record, error)
}
