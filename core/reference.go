package core

import (
	"fmt"
	"strings"
)

func ParseString(fullKey string) (Reference, error) {
	return Parse(strings.Split(fullKey, "/")...)
}

func Parse(parts ...string) (Reference, error) {
	if len(parts)%2 != 0 {
		return nil, fmt.Errorf("an uneven number of parts cannot result in a valid reference: %v", parts)
	}

	var sections []ReferenceSection
	for i := 0; i < len(parts); i += 2 {
		sections = append(sections, ReferenceSection{Kind: parts[i], Code: parts[i+1]})
	}

	return reference{sections}, nil
}

type ReferenceSection struct {
	Kind string
	Code string
}

type Reference interface {
	Parent() Reference
	Child(kind, code string) Reference
	Code() string
	Kind() string
	String() string
}

type reference struct {
	sections []ReferenceSection
}

func NewReference(kind, code string) Reference {
	return reference{[]ReferenceSection{{Kind: kind, Code: code}}}
}

func (k reference) Parent() Reference {
	if len(k.sections) == 0 {
		return nil
	}

	return reference{(k.sections)[:len(k.sections)-1]}
}

func (k reference) Child(kind, code string) Reference {
	return reference{append(k.sections, ReferenceSection{Kind: kind, Code: code})}
}

func (k reference) Code() string {
	if len(k.sections) == 0 {
		return ""
	}

	return k.sections[len(k.sections)-1].Code
}

func (k reference) Kind() string {
	if len(k.sections) == 0 {
		return ""
	}

	return k.sections[len(k.sections)-1].Kind
}

func (k reference) String() string {
	if len(k.sections) == 0 {
		return ""
	}

	var s string
	for _, section := range k.sections {
		s += section.Kind + "/" + section.Code + "/"
	}
	return s[:len(s)-1]
}

func (k reference) MarshalJSON() ([]byte, error) {
	return []byte(k.String()), nil
}

func (k reference) UnmarshalJSON(data []byte) error {
	parsed, err := Parse(string(data))
	if err != nil {
		return err
	}

	k.sections = parsed.(reference).sections
	return nil
}
