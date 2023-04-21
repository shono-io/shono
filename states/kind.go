package states

import (
	"strings"
)

func NewStateKind(concept, name string) Kind {
	return Kind{
		Concept: concept,
		Name:    name,
	}
}

func ParseStateKind(value string) *Kind {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return nil
	}

	return &Kind{
		Concept: parts[0],
		Name:    parts[1],
	}
}

type Kind struct {
	Concept string
	Name    string
}

func (sk Kind) String() string {
	return sk.Concept + ":" + sk.Name
}
