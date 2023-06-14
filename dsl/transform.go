package dsl

func Transform(language string, sourcecode string) TransformLogicStep {
	return TransformLogicStep{
		Mapping{
			Language:   language,
			Sourcecode: sourcecode,
		},
	}
}

type TransformLogicStep struct {
	Mapping
}

func (e TransformLogicStep) Kind() string {
	return "transform"
}
