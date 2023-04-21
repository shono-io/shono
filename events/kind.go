package events

import (
	"github.com/shono-io/go-shono/utils"
	"github.com/twmb/franz-go/pkg/kgo"
	"strings"
)

func EventKindFromHeader(headers []kgo.RecordHeader) *Kind {
	value := utils.Header(headers, KindHeader)
	if value == "" {
		return nil
	}

	return ParseEventKind(value)
}

func ParseEventKind(value string) *Kind {
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		return nil
	}

	return &Kind{
		Domain:  parts[0],
		Concept: parts[1],
		Name:    parts[2],
	}
}

func NewEventKind(domain, concept, name string) Kind {
	return Kind{
		Domain:  domain,
		Concept: concept,
		Name:    name,
	}
}

type Kind struct {
	Domain  string
	Concept string
	Name    string
}

func (ek Kind) String() string {
	return ek.Domain + ":" + ek.Concept + ":" + ek.Name
}
