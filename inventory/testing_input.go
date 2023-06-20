package inventory

import "github.com/shono-io/shono/commons"

func EventInput(eventRef commons.Reference, content map[string]any) TestInput {
	return TestInput{
		Metadata: map[string]string{
			"io_shono_kind": eventRef.String(),
		},
		Content: content,
	}
}

func RawInput(metadata map[string]string, content map[string]any) TestInput {
	return TestInput{
		Metadata: metadata,
		Content:  content,
	}
}
