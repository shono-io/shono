package commons

import (
	"fmt"
	"strings"
)

func Parse(parts ...string) (Key, error) {
	if len(parts)%2 != 0 {
		return nil, fmt.Errorf("an uneven number of parts cannot result in a valid key: %v", parts)
	}

	var sections []KeySection
	for i := 0; i < len(parts); i += 2 {
		sections = append(sections, KeySection{Kind: parts[i], Code: parts[i+1]})
	}

	return key{sections}, nil
}

func ParseString(fullKey string) (Key, error) {
	var result []string
	parts := strings.Split(fullKey, "__")
	for _, part := range parts {
		result = append(result, strings.Split(part, "_")...)
	}

	return Parse(result...)
}

type KeySection struct {
	Kind string
	Code string
}

type Key interface {
	Parent() Key
	Child(kind, code string) Key
	Code() string
	Kind() string
	String() string
	CodeString() string
}

type key struct {
	sections []KeySection
}

func NewKey(kind, code string) Key {
	return key{[]KeySection{{Kind: kind, Code: code}}}
}

func (k key) Parent() Key {
	if len(k.sections) == 0 {
		return nil
	}

	return key{(k.sections)[:len(k.sections)-1]}
}

func (k key) Child(kind, code string) Key {
	return key{append(k.sections, KeySection{Kind: kind, Code: code})}
}

func (k key) Code() string {
	if len(k.sections) == 0 {
		return ""
	}

	return k.sections[len(k.sections)-1].Code
}

func (k key) Kind() string {
	if len(k.sections) == 0 {
		return ""
	}

	return k.sections[len(k.sections)-1].Kind
}

func (k key) String() string {
	if len(k.sections) == 0 {
		return ""
	}

	var s string
	for _, section := range k.sections {
		s += section.Kind + "_" + section.Code + "__"
	}
	return s[:len(s)-2]
}

func (k key) CodeString() string {
	if len(k.sections) == 0 {
		return ""
	}

	var s string
	for _, section := range k.sections {
		s += section.Code + "__"
	}
	return s[:len(s)-2]
}

func (k key) MarshalJSON() ([]byte, error) {
	return []byte(k.String()), nil
}

func (k key) UnmarshalJSON(data []byte) error {
	parsed, err := Parse(string(data))
	if err != nil {
		return err
	}

	k.sections = parsed.(key).sections
	return nil
}
