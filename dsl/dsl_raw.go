package dsl

type RawLogic struct {
	Content string `yaml:"content"`
}

func (e RawLogic) Kind() string {
	return "raw"
}
