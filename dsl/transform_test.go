package dsl

import (
	"reflect"
	"testing"
)

func TestTransformLogicStep_MarshalBenthos(t *testing.T) {
	type fields struct {
		label   string
		mapping string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{
			"should fail if mapping is empty",
			fields{
				label:   "test",
				mapping: "",
			},
			nil, true,
		},
		{
			"should marshal with label",
			fields{
				label:   "test",
				mapping: "root = this",
			},
			map[string]any{
				"label":   "test",
				"mapping": "root = this",
			}, false,
		},
		{
			"should marshal without label",
			fields{
				mapping: "root = this",
			},
			map[string]any{
				"mapping": "root = this",
			}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := TransformLogicStep{
				label:   tt.fields.label,
				mapping: tt.fields.mapping,
			}
			got, err := e.MarshalBenthos("")
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
