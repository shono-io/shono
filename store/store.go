package store

import (
	"github.com/shono-io/shono"
)

type store struct {
	concept shono.Concept

	key         shono.Key
	name        string
	description string
}

func (e store) Key() shono.Key {
	return e.key
}

func (e store) Name() string {
	return e.name
}

func (e store) Description() string {
	return e.description
}

func newStoreOperation(id shono.StoreOpertationId, name, title, description string, isModifier bool, handler shono.StoreOperationHandler) shono.StoreOperation {
	return &storeOperation{
		id:          id,
		name:        name,
		title:       title,
		description: description,
		isModifier:  isModifier,
		handler:     handler,
	}
}

type storeOperation struct {
	id          shono.StoreOpertationId
	name        string
	title       string
	description string
	isModifier  bool
	handler     shono.StoreOperationHandler
}

func (s storeOperation) Id() shono.StoreOpertationId {
	return s.id
}

func (s storeOperation) Name() string {
	return s.name
}

func (s storeOperation) Title() string {
	return s.title
}

func (s storeOperation) Description() string {
	return s.description
}

func (s storeOperation) IsModifier() bool {
	return s.isModifier
}

func (s storeOperation) Handler() shono.StoreOperationHandler {
	return s.handler
}

func newResponse(code int, payload any) shono.StoreOperationResponse {
	return &response{
		code:    code,
		payload: payload,
	}
}

type response struct {
	code    int
	payload any
}

func (r response) Code() int {
	return r.code
}

func (r response) Payload() any {
	return r.payload
}
