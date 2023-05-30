package shono

import "strings"

type EventId string

func (e EventId) part(i int) string {
	parts := strings.Split(string(e), ":")
	if len(parts) != 3 {
		return ""
	}
	return parts[i]
}

func (e EventId) Scope() string {
	return e.part(0)
}

func (e EventId) Concept() string {
	return e.part(1)
}

func (e EventId) Code() string {
	return e.part(2)
}
