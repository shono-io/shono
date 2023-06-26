package dsl

import (
	"fmt"
	"github.com/shono-io/shono/commons"
	"strings"
)

func AsSuccessEvent(eventRef commons.Reference, status int, bloblangExpr string) TransformLogicStep {
	return TransformLogicStep{
		BloblangMapping(strings.TrimSpace(fmt.Sprintf(`
root = %s
meta shono_status = "%d"
meta shono_kind = %q
meta shono_backbone_topic = %q
`, bloblangExpr, status, eventRef.String(), eventRef.Parent().Parent().Code()))),
	}
}

func AsFailedEvent(eventRef commons.Reference, errorCode int, reason string) TransformLogicStep {
	return TransformLogicStep{
		BloblangMapping(strings.TrimSpace(fmt.Sprintf(`
root.reason = %s
meta shono_status = "%d"
meta shono_kind = %q
meta shono_backbone_topic = %q
`, reason, errorCode, eventRef.String(), eventRef.Parent().Parent().Code()))),
	}
}
