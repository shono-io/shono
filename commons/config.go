package commons

type PartialConfig struct {
	unmarshal func(any) error
}

func (pc *PartialConfig) UnmarshalYAML(unmarshal func(any) error) error {
	pc.unmarshal = unmarshal
	return nil
}

func (pc *PartialConfig) Unmarshal(v interface{}) error {
	return pc.unmarshal(v)
}
