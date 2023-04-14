package reaktor

import (
	sdk "github.com/shono-io/go-shono"
	"github.com/twmb/franz-go/pkg/kgo"
)

type StorageAccess string

var (
	StorageAccessRead         StorageAccess = "read"
	StorageAccessWrite        StorageAccess = "write"
	StorageAccessReadAndWrite StorageAccess = "read_and_write"
)

type HandlerFunc func(ctx Context, msg *kgo.Record) error

type Handler struct {
	HandlerInfo
	Kind sdk.EventKind
	Fn   HandlerFunc
}

func NewHandlerInfo(kind sdk.EventKind) *HandlerInfo {
	return &HandlerInfo{
		access: make(map[sdk.StateKind]StorageAccess),
	}
}

type HandlerInfo struct {
	emits  []sdk.EventKind
	access map[sdk.StateKind]StorageAccess
}

func (h *HandlerInfo) Emits(emits ...sdk.EventKind) *HandlerInfo {
	h.emits = append(h.emits, emits...)
	return h
}

func (h *HandlerInfo) Reads(stateKinds ...sdk.StateKind) *HandlerInfo {
	for _, stateKind := range stateKinds {
		h.access[stateKind] = StorageAccessRead
	}

	return h
}

func (h *HandlerInfo) Writes(stateKinds ...sdk.StateKind) *HandlerInfo {
	for _, stateKind := range stateKinds {
		h.access[stateKind] = StorageAccessWrite
	}
	return h
}

func (h *HandlerInfo) ReadsAndWrites(stateKinds ...sdk.StateKind) *HandlerInfo {
	for _, stateKind := range stateKinds {
		h.access[stateKind] = StorageAccessReadAndWrite
	}
	return h
}
