package graph

import "github.com/shono-io/shono/commons"

// == ENTITY ==========================================================================================================

type ReaktorOpt func(reaktor *Reaktor)

func WithName(name string) ReaktorOpt {
	return func(reaktor *Reaktor) {
		reaktor.name = name
	}
}

func WithReaktorDescription(description string) ReaktorOpt {
	return func(reaktor *Reaktor) {
		reaktor.description = description
	}
}

func WithOutputEvent(outputEventKeys ...commons.Key) ReaktorOpt {
	return func(reaktor *Reaktor) {
		reaktor.outputEventKeys = append(reaktor.outputEventKeys, outputEventKeys...)
	}
}

func WithStore(store ...Store) ReaktorOpt {
	return func(reaktor *Reaktor) {
		reaktor.stores = append(reaktor.stores, store...)
	}
}

func WithLogic(logic ...Logic) ReaktorOpt {
	return func(reaktor *Reaktor) {
		reaktor.logic = append(reaktor.logic, logic...)
	}
}

func WithTest(tests ...ReaktorTest) ReaktorOpt {
	return func(reaktor *Reaktor) {
		reaktor.tests = append(reaktor.tests, tests...)
	}
}

func NewReaktor(key commons.Key, inputEvent commons.Key, opts ...ReaktorOpt) Reaktor {
	result := Reaktor{
		key:           key,
		name:          key.Code(),
		inputEventKey: inputEvent,
	}

	for _, opt := range opts {
		opt(&result)
	}

	return result
}

type Reaktor struct {
	key             commons.Key
	name            string
	description     string
	inputEventKey   commons.Key
	outputEventKeys []commons.Key
	stores          []Store
	logic           []Logic
	tests           []ReaktorTest
}

func (r Reaktor) Key() commons.Key {
	return r.key
}

func (r Reaktor) Name() string {
	return r.name
}

func (r Reaktor) Description() string {
	return r.description
}

func (r Reaktor) InputEventKey() commons.Key {
	return r.inputEventKey
}

func (r Reaktor) OutputEventKeys() []commons.Key {
	return r.outputEventKeys
}

func (r Reaktor) Logic() []Logic {
	return r.logic
}

func (r Reaktor) Tests() []ReaktorTest {
	return r.tests
}

func (r Reaktor) Stores() []Store {
	return r.stores
}
