package dsl

import "strings"

func CleanMultilineStringWhitespace(sourcecode string) string {
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

	return sourcecode
}
