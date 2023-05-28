package backbone

import "github.com/shono-io/go-shono/event"

type RouteOpt func(r *Route)

func OverwrittenLog(log string) RouteOpt {
	return func(r *Route) {
		r.EventLog = log
	}
}

type Route struct {
	EventType *event.EventType
	EventLog  string
	Handler   EventHandler
}
