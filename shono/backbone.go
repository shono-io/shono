package shono

import (
	"fmt"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/memphis"
)

func NewBackbone(id string, kind string, props map[string]interface{}) (Backbone, error) {
	switch kind {
	case "memphis":
		return memphis.NewMemphisBackbone(id, props)
	}

	return nil, fmt.Errorf("unknown backbone kind: %s", kind)
}

type Backbone interface {
	go_shono.Writer

	Listen(r *go_shono.Router) error

	Apply(eid go_shono.EventId, event any) error

	Close()
}
