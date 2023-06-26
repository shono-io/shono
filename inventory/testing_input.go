package inventory

import (
	"fmt"
	"github.com/shono-io/shono/commons"
)

type TestInputOpt func(input *TestInput)

func WithTimestamp(timestamp int64) TestInputOpt {
	return func(input *TestInput) {
		input.Metadata["shono_timestamp"] = fmt.Sprintf("%d", timestamp)
	}
}

func WithEventRef(ref commons.Reference) TestInputOpt {
	return func(input *TestInput) {
		input.Metadata["shono_backbone_topic"] = ref.Parent().Parent().Code()
		input.Metadata["shono_kind"] = ref.String()
	}
}

func WithMetadata(key, value string) TestInputOpt {
	return func(input *TestInput) {
		input.Metadata[key] = value
	}
}

func NewTestInput(content map[string]any, opts ...TestInputOpt) TestInput {
	res := TestInput{
		Metadata: map[string]string{},
		Content:  content,
	}

	for _, opt := range opts {
		opt(&res)
	}

	return res
}

func RawInput(metadata map[string]string, content map[string]any) TestInput {
	return TestInput{
		Metadata: metadata,
		Content:  content,
	}
}
