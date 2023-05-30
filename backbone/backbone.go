package backbone

import (
	"context"
	"encoding/json"
	"github.com/shono-io/go-shono/event"
	"github.com/sirupsen/logrus"
	"time"
)

/*
EventHandler is a function that handles an event.

an error can be returned and is interpreted according to the following rules:
  - if the error is nil, the event is acked
  - if the error is a PoisonPillError, the event metadata is consulted to determine if the event should be acked or not
    and the application should panic or not
  - if the error is any other error, the application will panic and not commit the event
*/
type EventHandler func(ctx context.Context, evt *event.Event, logger *logrus.Logger) error

var NoopEventHandler = func(ctx context.Context, evt *event.Event, logger *logrus.Logger) error {
	b, _ := json.Marshal(evt)
	logger.Tracef("received event %s", b)
	return nil
}

type Backbone interface {
	MustWrite(ctx context.Context, evt *event.Event)

	Write(ctx context.Context, evt *event.Event) error

	On(eventType *event.EventType, fn EventHandler, opts ...RouteOpt) Backbone

	WaitFor(correlationId string, timeout time.Duration, possibleEvents ...*event.EventType) (*event.Event, error)

	Listen() error

	Close()
}

type Kind string

var (
	KafkaKind Kind = "kafka_franz"
)

type Config struct {
	Kind     Kind
	Settings map[string]any
}
