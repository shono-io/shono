package dsl

import (
	"fmt"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
)

func ListFromStore(scopeCode, conceptCode string, filters map[string]string) *StoreStepBuilder {
	return newStoreStepBuilder().
		Concept(scopeCode, conceptCode).
		Operation(StoreOperationList).
		Filters(filters)
}

func GetFromStore(scopeCode, conceptCode string, key string) *StoreStepBuilder {
	return newStoreStepBuilder().
		Concept(scopeCode, conceptCode).
		Operation(StoreOperationGet).
		Key(key)
}

func AddToStore(scopeCode, conceptCode string, key string) *StoreStepBuilder {
	return newStoreStepBuilder().
		Concept(scopeCode, conceptCode).
		Operation(StoreOperationAdd).
		Key(key)
}

func SetInStore(scopeCode, conceptCode string, key string) *StoreStepBuilder {
	return newStoreStepBuilder().
		Concept(scopeCode, conceptCode).
		Operation(StoreOperationSet).
		Key(key)
}

func RemoveFromStore(scopeCode, conceptCode string, key string) *StoreStepBuilder {
	return newStoreStepBuilder().
		Concept(scopeCode, conceptCode).
		Operation(StoreOperationDelete).
		Key(key)
}

func newStoreStepBuilder() *StoreStepBuilder {
	return &StoreStepBuilder{&StoreLogicStep{}}
}

type StoreStepBuilder struct {
	*StoreLogicStep
}

func (s *StoreStepBuilder) Operation(operation StoreOperation) *StoreStepBuilder {
	s.StoreLogicStep.Operation = operation
	return s
}

func (s *StoreStepBuilder) Concept(scopeCode, conceptCode string) *StoreStepBuilder {
	s.StoreLogicStep.Concept = inventory.NewConceptReference(scopeCode, conceptCode)
	return s
}

func (s *StoreStepBuilder) Label(label string) *StoreStepBuilder {
	s.StoreLogicStep.label = label
	return s
}

func (s *StoreStepBuilder) Filters(filters map[string]string) *StoreStepBuilder {
	s.StoreLogicStep.Filters = filters
	return s
}

func (s *StoreStepBuilder) Key(key string) *StoreStepBuilder {
	s.StoreLogicStep.Key = key
	return s
}

func (s *StoreStepBuilder) Build() inventory.LogicStep {
	return *s.StoreLogicStep
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
	label     string
	Concept   commons.Reference
	Operation StoreOperation
	Key       string
	Filters   map[string]string
}

func (s StoreLogicStep) Label() string {
	return s.label
}

func (s StoreLogicStep) MarshalBenthos(trace string) (map[string]any, error) {
	trace = fmt.Sprintf("%s/%s", trace, s.Kind())

	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", trace, err)
	}

	result := map[string]any{
		"concept":   s.Concept.String(),
		"operation": string(s.Operation),
	}

	if s.Key != "" {
		result["key"] = s.Key
	}

	if len(s.Filters) > 0 {
		result["filters"] = s.Filters
	}

	r := map[string]any{
		"store": result,
	}

	if s.label != "" {
		r["label"] = s.label
	}

	return r, nil
}

func (s StoreLogicStep) Kind() string {
	return "store"
}

func (s StoreLogicStep) Validate() error {
	if string(s.Operation) == "" {
		return fmt.Errorf("no operation defined")
	}

	if !s.Concept.IsValid() {
		return fmt.Errorf("no concept defined")
	}
	if s.Concept.Kind() != "concepts" {
		return fmt.Errorf("concept reference must point to a concept, not %s", s.Concept.Kind())
	}

	switch s.Operation {
	case StoreOperationList:
		if s.Key != "" {
			return fmt.Errorf("list operation does not accept a key")
		}
		if s.Filters == nil || len(s.Filters) == 0 {
			return fmt.Errorf("list operation requires filters")
		}
	case StoreOperationGet:
		if s.Key == "" {
			return fmt.Errorf("get operation requires a key")
		}
		if len(s.Filters) > 0 {
			return fmt.Errorf("get operation does not accept filters")
		}
	case StoreOperationAdd:
		if s.Key == "" {
			return fmt.Errorf("add operation requires a key")
		}
		if len(s.Filters) > 0 {
			return fmt.Errorf("add operation does not accept filters")
		}
	case StoreOperationSet:
		if s.Key == "" {
			return fmt.Errorf("set operation requires a key")
		}
		if len(s.Filters) > 0 {
			return fmt.Errorf("set operation does not accept filters")
		}
	case StoreOperationDelete:
		if s.Key == "" {
			return fmt.Errorf("delete operation requires a key")
		}
		if len(s.Filters) > 0 {
			return fmt.Errorf("delete operation does not accept filters")
		}
	}

	return nil
}
