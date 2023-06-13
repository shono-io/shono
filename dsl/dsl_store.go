package dsl

import (
	"github.com/shono-io/shono/graph"
)

func ListFromStore(scopeCode, conceptCode string, filters map[string]graph.Expression) StoreLogic {
	return StoreLogic{
		Concept: graph.ConceptReference{
			ScopeCode: scopeCode,
			Code:      conceptCode,
		},
		Operation: StoreOperationList,
		Filters:   filters,
	}
}

func GetFromStore(scopeCode, conceptCode string, key graph.Expression) StoreLogic {
	return StoreLogic{
		Concept: graph.ConceptReference{
			ScopeCode: scopeCode,
			Code:      conceptCode,
		},
		Operation: StoreOperationGet,
		Key:       &key,
	}
}

func AddToStore(scopeCode, conceptCode string, key graph.Expression, value ...Mapping) StoreLogic {
	return StoreLogic{
		Concept: graph.ConceptReference{
			ScopeCode: scopeCode,
			Code:      conceptCode,
		},
		Operation: StoreOperationAdd,
		Key:       &key,
		Value:     value,
	}
}

func SetInStore(scopeCode, conceptCode string, key graph.Expression, value ...Mapping) StoreLogic {
	return StoreLogic{
		Concept: graph.ConceptReference{
			ScopeCode: scopeCode,
			Code:      conceptCode,
		},
		Operation: StoreOperationSet,
		Key:       &key,
		Value:     value,
	}
}

func RemoveFromStore(scopeCode, conceptCode string, key graph.Expression) StoreLogic {
	return StoreLogic{
		Concept: graph.ConceptReference{
			ScopeCode: scopeCode,
			Code:      conceptCode,
		},
		Operation: StoreOperationDelete,
		Key:       &key,
	}
}

type StoreOperation string

var (
	StoreOperationList StoreOperation = "list"
	StoreOperationGet  StoreOperation = "get"

	StoreOperationAdd    StoreOperation = "add"
	StoreOperationSet    StoreOperation = "set"
	StoreOperationDelete StoreOperation = "delete"
)

type StoreLogic struct {
	Concept   graph.ConceptReference
	Operation StoreOperation
	Key       *graph.Expression
	Value     []Mapping
	Filters   map[string]graph.Expression
}

func (s StoreLogic) Kind() string {
	return "store"
}
