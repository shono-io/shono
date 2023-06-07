package graph

import (
	"github.com/shono-io/shono/commons"
)

type RequestType string

var (
	ListOperationType   RequestType = "list"
	GetOperationType    RequestType = "get"
	CreateOperationType RequestType = "create"
	UpdateOperationType RequestType = "update"
	DeleteOperationType RequestType = "delete"
)

type Request struct {
	Name          string
	Kind          RequestType
	Description   string
	RequestConfig *RequestConfig
}

type RequestConfig struct {
	StoreKey       commons.Key `json:"storeKey,omitempty"`
	PostEventKey   commons.Key `json:"postEventKey,omitempty"`
	PutEventKey    commons.Key `json:"putEventKey,omitempty"`
	DeleteEventKey commons.Key `json:"deleteEventKey,omitempty"`
}
