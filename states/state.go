package states

import (
	"encoding/json"
)

type State struct {
	kind      Kind
	valueType any
}

func (s *State) Kind() Kind {
	return s.kind
}

func (s *State) Encode(value any) ([]byte, error) {
	return json.Marshal(value)
}

func (s *State) Decode(b []byte, value any) error {
	return json.Unmarshal(b, value)
}

func (s *State) NewInstance() any {
	return s.valueType
}
