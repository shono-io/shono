package dsl

import (
	"context"
	"fmt"
	"github.com/shono-io/shono/graph"
)

func Transform(mappings ...Mapping) TransformLogic {
	return TransformLogic{
		Mappings: mappings,
	}
}

type TransformLogic struct {
	Mappings []Mapping
}

func (e TransformLogic) Kind() string {
	return "transform"
}

func MapRoot() Mapping {
	return Mapping{
		Field: "",
		Value: ToBloblang(`this`),
	}
}

func MapMeta(field graph.Expression, value Value) Mapping {
	return Mapping{
		Field: field,
		Value: value,
		Meta:  true,
	}
}

func Map(field graph.Expression, value Value) Mapping {
	return Mapping{
		Field: field,
		Value: value,
	}
}

type Mapping struct {
	Field graph.Expression `yaml:"field"`
	Value Value            `yaml:"value"`
	Meta  bool             `yaml:"meta"`
}

func (m Mapping) Generate(ctx context.Context) (string, error) {
	fieldName := string(m.Field)

	if m.Meta {
		fieldName = fmt.Sprintf("meta %s", fieldName)
	} else {
		if m.Field != "" {
			fieldName = fmt.Sprintf("root.%s", fieldName)
		} else {
			fieldName = "root"
		}
	}

	return fmt.Sprintf("%s = %s", fieldName, m.Value.Generate(ctx)), nil
}

type Value interface {
	Generate(ctx context.Context) string
}

func ToBloblang(value string) Value {
	return BlobValue{
		value: value,
	}
}

type BlobValue struct {
	value string
}

func (v BlobValue) Generate(ctx context.Context) string {
	return v.value
}

func AsConstant(value any) Value {
	switch value.(type) {
	case string:
		value = fmt.Sprintf("%q", value)
	}

	return ConstantValue{
		value: value,
	}
}

type ConstantValue struct {
	value any
}

func (v ConstantValue) Generate(ctx context.Context) string {
	return fmt.Sprintf("%v", v.value)
}

func AsEventReference(scopeCode, conceptCode, eventCode string) Value {
	return EventReferenceValue{
		reference: graph.EventReference{
			ScopeCode:   scopeCode,
			ConceptCode: conceptCode,
			Code:        eventCode,
		},
	}
}

type EventReferenceValue struct {
	reference graph.EventReference
}

func (v EventReferenceValue) Generate(ctx context.Context) string {
	return fmt.Sprintf("%q", v.reference)
}
