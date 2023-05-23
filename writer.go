package go_shono

import "context"

type Writer interface {
	MustWrite(ctx context.Context, correlationId string, evt *EventMeta, payload Payload)
	Write(ctx context.Context, correlationId string, evt *EventMeta, payload Payload) error
}
