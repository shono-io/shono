package shono

type rawBenthosConfigSection struct {
	yamlContent string
}

func (i *rawBenthosConfigSection) YAML() string {
	if i == nil {
		return ""
	}

	return i.yamlContent
}

func NewRawInput(rawYaml string) Input {
	return &RawInput{rawBenthosConfigSection{yamlContent: rawYaml}}
}

type RawInput struct {
	rawBenthosConfigSection
}

func NewRawOutput(rawYaml string) Output {
	return &RawOutput{rawBenthosConfigSection{yamlContent: rawYaml}}
}

type RawOutput struct {
	rawBenthosConfigSection
}

func NewRawProcessor(rawYaml string) Processor {
	return &RawProcessor{rawBenthosConfigSection{yamlContent: rawYaml}}
}

type RawProcessor struct {
	rawBenthosConfigSection
}
