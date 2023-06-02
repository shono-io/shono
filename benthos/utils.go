package benthos

import (
	"regexp"
)

var labelRegex = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func labelize(s string) string {
	return labelRegex.ReplaceAllString(s, "_")
}
