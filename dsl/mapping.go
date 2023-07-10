package dsl

import (
	"fmt"
	"strings"
)

func BloblangMapping(sourcecode string) Mapping {
	sourcecode = strings.TrimSpace(sourcecode)

	// -- cleanup bloblang newlines by taking the first line and checking how much whitespace there is. Then remove
	// -- that amount of whitespace from the beginning of each line.
	lines := strings.Split(sourcecode, "\n")
	if len(lines) > 1 {
		refLine := 0
		for len(lines[refLine]) == 0 {
			refLine++
			if refLine == len(lines) {
				break
			}
		}

		whitespace := ""
		for _, c := range lines[refLine] {
			if c == ' ' || c == '\t' {
				whitespace += string(c)
			} else {
				break
			}
		}

		for i := 1; i < len(lines); i++ {
			lines[i] = strings.TrimPrefix(lines[i], whitespace)
		}

		sourcecode = strings.Join(lines, "\n")
	}

	return Mapping{
		Language:   "bloblang",
		Sourcecode: sourcecode,
	}
}

type Mapping struct {
	Language   string
	Sourcecode string
}

func (m Mapping) Validate() error {
	if m.Language == "" {
		return fmt.Errorf("transform logic must have a language")
	}

	if m.Sourcecode == "" {
		return fmt.Errorf("transform logic must have sourcecode")
	}

	return nil
}
