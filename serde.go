package go_shono

import "encoding/json"

const CorrelationHeader = "io.shono.correlation"
const KindHeader = "io.shono.kind"

type Serde interface {
	Encode(payload any) ([]byte, error)
	Decode(value []byte, target any) error
}

func JsonSerde() Serde {
	return jsonSerde{}
}

type jsonSerde struct {
}

func (j jsonSerde) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (j jsonSerde) Decode(bytes []byte, i interface{}) error {
	return json.Unmarshal(bytes, i)
}
