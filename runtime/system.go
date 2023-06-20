package runtime

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
)

func NewSystemReference(code string) commons.Reference {
	return commons.NewReference("systems", code)
}

type InputSystem interface {
	System
	AsInput(config map[string]any) (map[string]any, error)
}

type OutputSystem interface {
	System
	AsOutput(config map[string]any) (map[string]any, error)
}

type System interface {
	inventory.Node
}
