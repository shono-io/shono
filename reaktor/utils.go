package reaktor

import (
	goshono "github.com/shono-io/go-shono"
	"github.com/twmb/franz-go/pkg/kgo"
)

func Header(headers []kgo.RecordHeader, key string) string {
	for _, header := range headers {
		if header.Key == key {
			return string(header.Value)
		}
	}

	return ""
}

func EventKindFromHeader(headers []kgo.RecordHeader) *goshono.EventKind {
	value := Header(headers, goshono.KindHeader)
	if value == "" {
		return nil
	}

	return goshono.ParseEventKind(value)
}

func StateKindFromHeader(headers []kgo.RecordHeader) *goshono.StateKind {
	value := Header(headers, goshono.KindHeader)
	if value == "" {
		return nil
	}

	return goshono.ParseStateKind(value)
}
