package dsl

import (
	"fmt"
	"github.com/shono-io/shono/commons"
)

func ListFromStore(scopeCode, conceptCode string, filters map[string]string) StoreLogicStep {
	return StoreLogicStep{
		Concept:   commons.NewReference("scopes", scopeCode).Child("concepts", conceptCode),
		Operation: StoreOperationList,
		Filters:   filters,
	}
}

func GetFromStore(scopeCode, conceptCode string, key string) StoreLogicStep {
	return StoreLogicStep{
		Concept:   commons.NewReference("scopes", scopeCode).Child("concepts", conceptCode),
		Operation: StoreOperationGet,
		Key:       key,
	}
}

func AddToStore(scopeCode, conceptCode string, key string, value Mapping) StoreLogicStep {
	return StoreLogicStep{
		Concept:   commons.NewReference("scopes", scopeCode).Child("concepts", conceptCode),
		Operation: StoreOperationAdd,
		Key:       key,
		Value:     &value,
	}
}

func SetInStore(scopeCode, conceptCode string, key string, value Mapping) StoreLogicStep {
	return StoreLogicStep{
		Concept:   commons.NewReference("scopes", scopeCode).Child("concepts", conceptCode),
		Operation: StoreOperationSet,
		Key:       key,
		Value:     &value,
	}
}

func RemoveFromStore(scopeCode, conceptCode string, key string) StoreLogicStep {
	return StoreLogicStep{
		Concept:   commons.NewReference("scopes", scopeCode).Child("concepts", conceptCode),
		Operation: StoreOperationDelete,
		Key:       key,
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

type StoreLogicStep struct {
	Concept   commons.Reference
	Operation StoreOperation
	Key       string
	Value     *Mapping
	Filters   map[string]string
}

func (s StoreLogicStep) Kind() string {
	return "store"
}

func (s StoreLogicStep) Validate() error {
	if s.Operation == "" {
		return fmt.Errorf("no operation defined")
	}

	if s.Concept == nil {
		return fmt.Errorf("no concept defined")
	}

	switch s.Operation {
	case StoreOperationList:
		if s.Key != "" {
			return fmt.Errorf("list operation does not accept a key")
		}

		if s.Value != nil {
			return fmt.Errorf("list operation does not accept a value")
		}
	case StoreOperationGet:
		if s.Key == "" {
			return fmt.Errorf("get operation requires a key")
		}

		if s.Value != nil {
			return fmt.Errorf("get operation does not accept a value")
		}

		if len(s.Filters) > 0 {
			return fmt.Errorf("get operation does not accept filters")
		}
	case StoreOperationAdd:
		if s.Key == "" {
			return fmt.Errorf("add operation requires a key")
		}
		if s.Value == nil {
			return fmt.Errorf("add operation requires a value")
		}
		if len(s.Filters) > 0 {
			return fmt.Errorf("add operation does not accept filters")
		}
	case StoreOperationSet:
		if s.Key == "" {
			return fmt.Errorf("set operation requires a key")
		}
		if s.Value == nil {
			return fmt.Errorf("set operation requires a value")
		}
		if len(s.Filters) > 0 {
			return fmt.Errorf("set operation does not accept filters")
		}
	case StoreOperationDelete:
		if s.Key == "" {
			return fmt.Errorf("delete operation requires a key")
		}
		if s.Value != nil {
			return fmt.Errorf("delete operation does not accept a value")
		}
		if len(s.Filters) > 0 {
			return fmt.Errorf("delete operation does not accept filters")
		}
	}

	return nil
}
