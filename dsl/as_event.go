package dsl

import (
	"fmt"
	"github.com/shono-io/shono/commons"
)

func AsEvent(eventRef commons.Reference) TransformLogicStep {
	return TransformLogicStep{
		BloblangMapping(fmt.Sprintf(`meta io_shono_kind = %q
meta backbone_topic = %q
root = this
`, eventRef.String(), eventRef.Parent().Parent().Code())),
	}
}
