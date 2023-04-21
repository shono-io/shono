package utils

import "github.com/twmb/franz-go/pkg/kgo"

func Header(headers []kgo.RecordHeader, key string) string {
	for _, header := range headers {
		if header.Key == key {
			return string(header.Value)
		}
	}

	return ""
}
