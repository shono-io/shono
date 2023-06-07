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

type RequestOpt func(*Request)

func WithRequestDescription(description string) RequestOpt {
	return func(r *Request) {
		r.Description = description
	}
}

func WithRequestName(name string) RequestOpt {
	return func(r *Request) {
		r.Name = name
	}
}

func NewGetRequest(storeKey commons.Key, opts ...RequestOpt) Request {
	return newRequest(GetOperationType, storeKey, nil, opts...)
}

func NewListRequest(storeKey commons.Key, opts ...RequestOpt) Request {
	return newRequest(ListOperationType, storeKey, nil, opts...)
}

func NewAddRequest(eventKey commons.Key, opts ...RequestOpt) Request {
	return newRequest(CreateOperationType, nil, eventKey, opts...)
}

func NewSetRequest(eventKey commons.Key, opts ...RequestOpt) Request {
	return newRequest(UpdateOperationType, nil, eventKey, opts...)
}

func NewDeleteRequest(eventKey commons.Key, opts ...RequestOpt) Request {
	return newRequest(DeleteOperationType, nil, eventKey, opts...)
}

func newRequest(kind RequestType, storeKey commons.Key, eventKey commons.Key, opts ...RequestOpt) Request {
	result := Request{
		Kind:     kind,
		StoreKey: storeKey,
		EventKey: eventKey,
	}

	for _, opt := range opts {
		opt(&result)
	}

	return result
}

type Request struct {
	Kind        RequestType
	Name        string
	Description string

	StoreKey commons.Key
	EventKey commons.Key
}

type GetRequest struct {
	StoreKey commons.Key
}

type ListRequest struct {
	StoreKey commons.Key
}

type AddRequest struct {
	EventKey commons.Key
}

type SetRequest struct {
	EventKey commons.Key
}

type DeleteRequest struct {
	EventKey commons.Key
}
