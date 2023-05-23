package kafka

import (
	"context"
	"fmt"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/utils"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

type CallbackNotification struct {
	Event    go_shono.EventId
	Value    *any
	timedOut bool
}

type callback struct {
	c      chan CallbackNotification
	events []go_shono.EventId
	t      *time.Timer
	cancel chan any
}

func (c *callback) matches(id go_shono.EventId) bool {
	for _, eid := range c.events {
		if eid == id {
			return true
		}
	}

	return false
}

func (c *callback) resolve(eid go_shono.EventId, evt any) {
	c.c <- CallbackNotification{
		Event: eid,
		Value: &evt,
	}
	c.t.Stop()
	close(c.cancel)
}

func NewRunner(name string, r *go_shono.Router, opts ...kgo.Opt) *Runner {
	opts = append(opts, kgo.ConsumerGroup(name), kgo.DisableAutoCommit(), kgo.ConsumeTopics(r.Scopes()...))

	return &Runner{
		name:      name,
		opts:      opts,
		r:         r,
		callbacks: make(map[string]callback),
	}
}

type Runner struct {
	name string
	opts []kgo.Opt

	r *go_shono.Router

	callbacks map[string]callback

	kc     *kgo.Client
	Writer *Writer
}

func (r *Runner) Close() {
	if r.kc != nil {
		r.kc.Close()
	}
}

func (r *Runner) handleRecord(ctx context.Context, msg *kafkaMsg) error {
	if msg.Empty() {
		return nil
	}

	// -- check if there is a kind header
	eid := utils.Header(msg.Record.Headers, go_shono.KindHeader)
	if eid == "" {
		return fmt.Errorf("no event kind header found")
	}

	corId := utils.Header(msg.Record.Headers, go_shono.CorrelationHeader)
	if corId == "" {
		return fmt.Errorf("no correlation id header found")
	}

	pctx := go_shono.WithCorrelationId(ctx, corId)
	pctx = go_shono.WithKey(pctx, string(msg.Record.Key))

	r.r.Process(pctx, go_shono.EventId(eid), msg.Record.Value)

	msg.Ack()

	// -- check if there is a callback for the correlation id
	if cb, fnd := r.callbacks[corId]; fnd {
		if cb.matches(go_shono.EventId(eid)) {
			// -- decode the message
			res, _, err := r.r.Decode(go_shono.EventId(eid), msg.Record.Value)
			if err != nil {
				return fmt.Errorf("failed to resolve callback: failed to decode message: %v", err)
			}

			cb.resolve(go_shono.EventId(eid), res)
		}
	}

	return nil
}

func (r *Runner) RegisterCallback(correlationId string, timeout time.Duration, events ...*go_shono.EventMeta) (chan CallbackNotification, error) {
	evts := make([]go_shono.EventId, len(events))
	for i, evt := range events {
		evts[i] = evt.EventId
	}

	// -- check if the callback already exists
	if _, fnd := r.callbacks[correlationId]; fnd {
		return nil, fmt.Errorf("callback already exists for correlation id %s", correlationId)
	}

	// -- create the channel
	ch := make(chan CallbackNotification)
	r.callbacks[correlationId] = callback{
		c:      ch,
		events: evts,
		t:      time.NewTimer(timeout),
		cancel: make(chan any),
	}

	// -- start a timer to remove the callback on timeout
	go func(cid string) {
		cb, fnd := r.callbacks[cid]
		if !fnd {
			return
		}

		select {
		case <-cb.cancel:
		case <-cb.t.C:
			cb.c <- CallbackNotification{
				timedOut: true,
			}
		}

		delete(r.callbacks, cid)
	}(correlationId)

	return ch, nil
}

func (r *Runner) Run() error {
	if r.kc != nil {
		return fmt.Errorf("runner already running")
	}

	// -- create the kgo client
	kc, err := kgo.NewClient(r.opts...)
	if err != nil {
		return fmt.Errorf("failed to create kafka client: %Writer", err)
	}
	r.kc = kc

	r.Writer = &Writer{kc: kc}

	ctx := context.Background()

	go func() {
		for {
			fetches := kc.PollFetches(ctx)
			if fetches.IsClientClosed() {
				return
			}

			// -- panic if we were unable to fetch messages
			if err := fetches.Err(); err != nil {
				panic(fmt.Errorf("failed to fetch: %w", err))
			}

			for _, record := range fetches.Records() {
				m := &kafkaMsg{
					kc:     kc,
					Record: record,
				}

				err := r.handleRecord(ctx, m)
				if err != nil {
					panic(fmt.Errorf("failed to handle record: %Writer", err))
				}
			}
		}
	}()

	return nil
}

type kafkaMsg struct {
	kc     *kgo.Client
	Record *kgo.Record
}

func (m *kafkaMsg) Ack() {
	if err := m.kc.CommitRecords(context.Background(), m.Record); err != nil {
		panic(fmt.Errorf("failed to ack message: %Writer", err))
	}
}

func (m *kafkaMsg) Empty() bool {
	return m.Record == nil || m.Record.Value == nil
}
