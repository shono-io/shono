package dsl

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
	"reflect"
	"testing"
)

func TestStoreLogicStep_MarshalBenthos(t *testing.T) {
	type fields struct {
		label     string
		Concept   commons.Reference
		Operation StoreOperation
		Key       string
		Filters   map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{
			"should fail if concept reference is invalid",
			fields{
				Concept:   commons.Reference(""),
				Operation: StoreOperationGet,
				Key:       "test",
			},
			nil, true,
		},
		{
			"should fail if the concept reference is actually pointing to a scope",
			fields{
				Concept:   inventory.NewScopeReference("s"),
				Operation: StoreOperationGet,
				Key:       "test",
			},
			nil, true,
		},
		{
			"should fail if the concept reference is actually pointing to an event",
			fields{
				Concept:   inventory.NewEventReference("s", "c", "e"),
				Operation: StoreOperationGet,
				Key:       "test",
			},
			nil, true,
		},
		{
			"should fail if an operation is not specified",
			fields{
				Concept: inventory.NewConceptReference("s", "c"),
				Key:     "test",
			},
			nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StoreLogicStep{
				label:     tt.fields.label,
				Concept:   tt.fields.Concept,
				Operation: tt.fields.Operation,
				Key:       tt.fields.Key,
				Filters:   tt.fields.Filters,
			}
			got, err := s.MarshalBenthos("")
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBenthos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalBenthos() got = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("LIST", testListOperation)
	t.Run("GET", testGetOperation)
	t.Run("ADD", testAddOperation)
	t.Run("SET", testSetOperation)
	t.Run("DELETE", testDeleteOperation)
}

func testListOperation(t *testing.T) {
	op := StoreOperationList

	type fields struct {
		label     string
		Concept   commons.Reference
		Operation StoreOperation
		Key       string
		Filters   map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{
			"should succeed if a valid operation is specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
				Filters: map[string]string{
					"a": "b",
				},
			},
			map[string]any{
				"store": map[string]any{
					"concept":   "scopes/s/concepts/c",
					"operation": string(op),
					"filters":   map[string]string{"a": "b"},
				},
			}, false,
		},
		{
			"should fail if no filters are specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
			},
			nil, true,
		},
		{
			"should fail if a key was specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
				Key:       "test",
				Filters:   map[string]string{"a": "b"},
			},
			nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StoreLogicStep{
				label:     tt.fields.label,
				Concept:   tt.fields.Concept,
				Operation: tt.fields.Operation,
				Key:       tt.fields.Key,
				Filters:   tt.fields.Filters,
			}
			got, err := s.MarshalBenthos("")
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBenthos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalBenthos() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func testGetOperation(t *testing.T) {
	op := StoreOperationGet

	type fields struct {
		label     string
		Concept   commons.Reference
		Operation StoreOperation
		Key       string
		Filters   map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{
			"should succeed if a valid operation is specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
				Key:       "test",
			},
			map[string]any{
				"store": map[string]any{
					"concept":   "scopes/s/concepts/c",
					"operation": string(op),
					"key":       "test",
				},
			}, false,
		},
		{
			"should fail if no key is specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
			},
			nil, true,
		},
		{
			"should fail if filters are specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
				Key:       "test",
				Filters:   map[string]string{"a": "b"},
			},
			nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StoreLogicStep{
				label:     tt.fields.label,
				Concept:   tt.fields.Concept,
				Operation: tt.fields.Operation,
				Key:       tt.fields.Key,
				Filters:   tt.fields.Filters,
			}
			got, err := s.MarshalBenthos("")
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBenthos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalBenthos() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func testAddOperation(t *testing.T) {
	op := StoreOperationAdd

	type fields struct {
		label     string
		Concept   commons.Reference
		Operation StoreOperation
		Key       string
		Filters   map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{
			"should succeed if a valid operation is specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
				Key:       "test",
			},
			map[string]any{
				"store": map[string]any{
					"concept":   "scopes/s/concepts/c",
					"operation": string(op),
					"key":       "test",
				},
			}, false,
		},
		{
			"should fail if no key is specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
			},
			nil, true,
		},
		{
			"should fail if filters are specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
				Key:       "test",
				Filters:   map[string]string{"a": "b"},
			},
			nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StoreLogicStep{
				label:     tt.fields.label,
				Concept:   tt.fields.Concept,
				Operation: tt.fields.Operation,
				Key:       tt.fields.Key,
				Filters:   tt.fields.Filters,
			}
			got, err := s.MarshalBenthos("")
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBenthos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalBenthos() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func testSetOperation(t *testing.T) {
	op := StoreOperationSet

	type fields struct {
		label     string
		Concept   commons.Reference
		Operation StoreOperation
		Key       string
		Filters   map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{
			"should succeed if a valid operation is specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
				Key:       "test",
			},
			map[string]any{
				"store": map[string]any{
					"concept":   "scopes/s/concepts/c",
					"operation": string(op),
					"key":       "test",
				},
			}, false,
		},
		{
			"should fail if no key is specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
			},
			nil, true,
		},
		{
			"should fail if filters are specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
				Key:       "test",
				Filters:   map[string]string{"a": "b"},
			},
			nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StoreLogicStep{
				label:     tt.fields.label,
				Concept:   tt.fields.Concept,
				Operation: tt.fields.Operation,
				Key:       tt.fields.Key,
				Filters:   tt.fields.Filters,
			}
			got, err := s.MarshalBenthos("")
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBenthos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalBenthos() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func testDeleteOperation(t *testing.T) {
	op := StoreOperationDelete

	type fields struct {
		label     string
		Concept   commons.Reference
		Operation StoreOperation
		Key       string
		Filters   map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{
			"should succeed if a valid operation is specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
				Key:       "test",
			},
			map[string]any{
				"store": map[string]any{
					"concept":   "scopes/s/concepts/c",
					"operation": string(op),
					"key":       "test",
				},
			}, false,
		},
		{
			"should fail if no key is specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
			},
			nil, true,
		},
		{
			"should fail if filters are specified",
			fields{
				Concept:   inventory.NewConceptReference("s", "c"),
				Operation: op,
				Key:       "test",
				Filters:   map[string]string{"a": "b"},
			},
			nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StoreLogicStep{
				label:     tt.fields.label,
				Concept:   tt.fields.Concept,
				Operation: tt.fields.Operation,
				Key:       tt.fields.Key,
				Filters:   tt.fields.Filters,
			}
			got, err := s.MarshalBenthos("")
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBenthos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalBenthos() got = %v, want %v", got, tt.want)
			}
		})
	}
}
