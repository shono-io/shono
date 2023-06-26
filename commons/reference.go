package commons

import (
	"fmt"
	"strings"
)

var ReferenceSeparator = "/"

func ParseString(fullKey string) (Reference, error) {
	if (fullKey == "") || (fullKey == ReferenceSeparator) {
		return "", fmt.Errorf("empty Reference")
	}

	return Parse(strings.Split(strings.TrimPrefix(fullKey, ReferenceSeparator), ReferenceSeparator)...)
}

func Parse(parts ...string) (Reference, error) {
	if len(parts)%2 != 0 {
		return "", fmt.Errorf("an uneven number of parts cannot result in a valid Reference: %v", parts)
	}

	return Reference(strings.Join(parts, ReferenceSeparator)), nil
}

type Reference string

func NewReference(kind, code string) Reference {
	return Reference(fmt.Sprintf("%s%s%s", kind, ReferenceSeparator, code))
}

func (k Reference) parts() []string {
	return strings.Split(string(k), ReferenceSeparator)
}

func (k Reference) IsValid() bool {
	return len(k) > 0 && len(k.parts())%2 == 0
}

func (k Reference) Parent() Reference {
	p := k.parts()
	if len(p) <= 2 {
		return ""
	}

	return Reference(strings.Join(p[:len(p)-2], ReferenceSeparator))
}

func (k Reference) Child(kind, code string) Reference {
	return Reference(fmt.Sprintf("%s%s%s%s%s", k, ReferenceSeparator, kind, ReferenceSeparator, code))
}

func (k Reference) Code() string {
	p := k.parts()
	if len(p) == 0 {
		return ""
	}

	return p[len(p)-1]
}

func (k Reference) Kind() string {
	p := k.parts()
	if len(p) < 2 {
		return ""
	}

	return p[len(p)-2]
}

func (k Reference) String() string {
	return string(k)
}
