package dsl

func Transform(mapping Mapping) TransformLogicStep {
	return TransformLogicStep{
		mapping,
	}
}

type TransformLogicStep struct {
	Mapping
}

func (e TransformLogicStep) Kind() string {
	return "transform"
}
