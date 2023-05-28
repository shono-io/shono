package event

import (
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/utils"
	"github.com/sirupsen/logrus"
	"time"
)

type Opt func(*Event)

func WithCorrelationIdFromEvent(evt *Event) Opt {
	return func(e *Event) {
		if corId, fnd := evt.metadata[go_shono.CorrelationHeader]; fnd {
			e.metadata[go_shono.CorrelationHeader] = corId
		}
	}
}

func WithCorrelationId(corId string) Opt {
	return func(e *Event) {
		e.metadata[go_shono.CorrelationHeader] = corId
	}
}

func WithAcker(acker func()) Opt {
	return func(e *Event) {
		e.acker = acker
	}
}

func WithNacker(nacker func(reason error, continueProcessing bool)) Opt {
	return func(e *Event) {
		e.nacker = nacker
	}
}

func WithMetadata(metadata map[string]any) Opt {
	return func(e *Event) {
		for k, v := range metadata {
			e.metadata[k] = v
		}
	}
}

func WithMetadataField(key string, value any) Opt {
	return func(e *Event) {
		e.metadata[key] = value
	}
}

func WithHeaderFromEvent(evt *Event, header string) Opt {
	return func(e *Event) {
		if val, fnd := evt.metadata[header]; fnd {
			e.metadata[header] = val
		}
	}
}

func NewEvent(et *EventType, value any, opts ...Opt) *Event {
	v, ok := value.([]byte)
	if !ok {
		v = utils.MustReturn(et.Encode(value))
	}

	res := &Event{
		EventId:  &et.EventId,
		t:        et,
		value:    v,
		metadata: map[string]any{},
	}

	for _, opt := range opts {
		opt(res)
	}

	// -- write the correct event id to the kind header
	res.metadata[go_shono.KindHeader] = string(*res.EventId)

	return res
}

type Event struct {
	*EventId
	t        *EventType
	value    []byte
	metadata map[string]any
	acker    func()
	nacker   func(reason error, redeliver bool)
}

func (e *Event) Is(t *EventType) bool {
	return e.t == t
}

func (e *Event) Raw() []byte {
	return e.value
}

func (e *Event) Value(result any) error {
	return e.t.Decode(e.value, &result)
}

func (e *Event) Meta(key string) any {
	return e.metadata[key]
}

func (e *Event) MetaString(key string) string {
	if e.metadata[key] == nil {
		return ""
	}

	return e.metadata[key].(string)
}

func (e *Event) MetaTime(key string) *time.Time {
	if e.metadata[key] == nil {
		return nil
	}

	return e.metadata[key].(*time.Time)
}

func (e *Event) MetaKeys() []string {
	keys := make([]string, 0, len(e.metadata))
	for k := range e.metadata {
		keys = append(keys, k)
	}

	return keys
}

func (e *Event) Ack() {
	e.acker()
}

func (e *Event) Skip(reason error) {
	e.acker()
	logrus.Warnf("skipping %s event with key %q: %v", *e.EventId, e.MetaString("kafka_key"), reason)
}

func (e *Event) Redeliver(reason error) {
	e.nacker(reason, true)
}

func (e *Event) Panic(reason error) {
	e.nacker(reason, false)
}
