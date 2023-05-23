package kafka

import (
	"context"
	"fmt"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/utils"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

func NewBackbone(id string, opts ...kgo.Opt) *KafkaBackbone {
	return &KafkaBackbone{
		id:   id,
		opts: opts,
	}
}

type KafkaBackbone struct {
	id   string
	opts []kgo.Opt

	run *Runner
}

func (k *KafkaBackbone) MustWrite(ctx context.Context, correlationId string, evt *go_shono.EventMeta, payload go_shono.Payload) {
	if k.run.Writer == nil {
		panic(fmt.Errorf("backbone not initialized"))
	}

	k.run.Writer.MustWrite(ctx, correlationId, evt, payload)
}

func (k *KafkaBackbone) Write(ctx context.Context, correlationId string, evt *go_shono.EventMeta, payload go_shono.Payload) error {
	if k.run.Writer == nil {
		panic(fmt.Errorf("backbone not initialized"))
	}

	return k.run.Writer.Write(ctx, correlationId, evt, payload)
}

func (k *KafkaBackbone) Listen(r *go_shono.Router) error {
	if k.run != nil {
		return fmt.Errorf("already listening")
	}

	k.run = NewRunner(k.id, r, k.opts...)
	return k.run.Run()
}

func (k *KafkaBackbone) WaitFor(correlationId string, timeout time.Duration, possibleEvents ...*go_shono.EventMeta) (go_shono.EventId, any, error) {
	c, err := k.run.RegisterCallback(correlationId, timeout, possibleEvents...)
	if err != nil {
		return "", nil, err
	}

	// -- wait for the channel to complete
	res := <-c
	if res.timedOut {
		return "", nil, utils.ErrTimeout
	} else {
		return res.Event, *res.Value, nil
	}
}

func (k *KafkaBackbone) Apply(eid go_shono.EventId, event any) error {
	//switch eid {
	//case events.ScopeCreated.EventId:
	//	//return b.onScopeCreated(event.(*events.ScopeCreatedEvent))
	//case events.ScopeDeleted.EventId:
	//	//return b.onScopeDeleted(event.(*events.ScopeDeletedEvent))
	//}

	return nil
}

func (k *KafkaBackbone) Close() {
	if k.run != nil {
		k.run.Close()
	}
}
