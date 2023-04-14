package go_shono

import "strings"

func ParseEventKind(value string) *EventKind {
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		return nil
	}

	return &EventKind{
		Domain:  parts[0],
		Concept: parts[1],
		Name:    parts[2],
	}
}

type EventKind struct {
	Domain  string
	Concept string
	Name    string
}

func (ek EventKind) String() string {
	return ek.Domain + ":" + ek.Concept + ":" + ek.Name
}

func ParseStateKind(value string) *StateKind {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return nil
	}

	return &StateKind{
		Concept: parts[0],
		Name:    parts[1],
	}
}

type StateKind struct {
	Concept string
	Name    string
}

func (sk StateKind) String() string {
	return sk.Concept + ":" + sk.Name
}
