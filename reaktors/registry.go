package reaktors

import "github.com/shono-io/go-shono/events"

type ReaktorRegistry interface {
	ReaktorFor(kind events.Kind) (*ReaktorInfo, error)
}
