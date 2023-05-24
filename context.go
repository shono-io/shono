package go_shono

import (
	"context"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

type SendOpt func(r *kgo.Record)

type Context interface {
	context.Context

	Failed(evt EventMeta, err error)
	Send(evt EventMeta, payload any, opts ...SendOpt)

	Resource(id string) any
}

func WithOrganization(ctx context.Context, org string) context.Context {
	return context.WithValue(ctx, "organization", org)
}

func OrganizationFromContext(ctx context.Context) string {
	res := ctx.Value("organization")
	if res == nil {
		return ""
	}

	return res.(string)
}

func WithKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, "key", &key)
}

func KeyFromContext(ctx context.Context) *string {
	res := ctx.Value("key")
	if res == nil {
		return nil
	}

	return res.(*string)
}

func WithCorrelationId(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, "correlation_id", &cid)
}

func CorrelationIdFromContext(ctx context.Context) *string {
	res := ctx.Value("correlation_id")
	if res == nil {
		return nil
	}

	return res.(*string)
}

func WithEvent(ctx context.Context, evt *EventMeta) context.Context {
	return context.WithValue(ctx, "event", &evt)
}

func EventFromContext(ctx context.Context) *EventMeta {
	res := ctx.Value("event")
	if res == nil {
		return nil
	}

	return res.(*EventMeta)
}

func WithEventTimestamp(ctx context.Context, time time.Time) context.Context {
	return context.WithValue(ctx, "event_timestamp", &time)
}

func EventTimestampFromContext(ctx context.Context) *time.Time {
	res := ctx.Value("event_timestamp")
	if res == nil {
		return nil
	}

	return res.(*time.Time)
}
