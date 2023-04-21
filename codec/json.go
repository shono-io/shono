package codec

import "encoding/json"

type Json struct {
}

func (j Json) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (j Json) Decode(bytes []byte, i interface{}) error {
	return json.Unmarshal(bytes, i)
}
