package graph

import (
	"github.com/shono-io/shono/core"
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

func NewGetRequest(storeKey core.Reference, opts ...RequestOpt) Request {
	return newRequest(GetOperationType, storeKey, nil, opts...)
}

func NewListRequest(storeKey core.Reference, opts ...RequestOpt) Request {
	return newRequest(ListOperationType, storeKey, nil, opts...)
}

func NewAddRequest(eventKey core.Reference, opts ...RequestOpt) Request {
	return newRequest(CreateOperationType, nil, eventKey, opts...)
}

func NewSetRequest(eventKey core.Reference, opts ...RequestOpt) Request {
	return newRequest(UpdateOperationType, nil, eventKey, opts...)
}

func NewDeleteRequest(eventKey core.Reference, opts ...RequestOpt) Request {
	return newRequest(DeleteOperationType, nil, eventKey, opts...)
}

func newRequest(kind RequestType, storeKey core.Reference, eventKey core.Reference, opts ...RequestOpt) Request {
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

	StoreKey core.Reference
	EventKey core.Reference
}

type GetRequest struct {
	StoreKey core.Reference
}

type ListRequest struct {
	StoreKey core.Reference
}

type AddRequest struct {
	EventKey core.Reference
}

type SetRequest struct {
	EventKey core.Reference
}

type DeleteRequest struct {
	EventKey core.Reference
}
