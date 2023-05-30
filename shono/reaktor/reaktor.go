package reaktor

import (
	"fmt"
	"github.com/shono-io/go-shono/shono"
	"github.com/shono-io/go-shono/shono/benthos"
)

func NewReaktor(scopeCode string, code string, name string, description string, backbone shono.Backbone, inputEvents []shono.EventId, outputEvents []shono.EventId, processors []benthos.Part) shono.Reaktor {
	var opts []benthos.ReaktorOpt
	opts = append(opts, benthos.WithProcessor(processors...))

	input := NewShonoInput(fmt.Sprintf("%s-%s-reaktor", scopeCode, code), backbone, inputEvents...)
	opts = append(opts, benthos.WithInput(input))

	output := NewOutput(backbone, outputEvents...)
	opts = append(opts, benthos.WithOutput(output))

	return benthos.NewBenthosReaktor(scopeCode, code, name, description, opts...)
}
