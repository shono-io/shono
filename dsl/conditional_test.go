package dsl

import (
	"github.com/shono-io/shono/inventory"
	"reflect"
	"testing"
)

func TestConditionalLogicStep_MarshalBenthos(t *testing.T) {
	type fields struct {
		label       string
		Cases       []ConditionalCase
		DefaultCase *ConditionalCase
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{
			"should fail if no cases are provided",
			fields{
				"my_label",
				[]ConditionalCase{},
				nil,
			},
			nil,
			true,
		},
		{
			"should fail if no processors are provided but a default case is",
			fields{
				"my_label",
				[]ConditionalCase{},
				&ConditionalCase{
					"", []inventory.LogicStep{
						Log(LogLevelInfo).Message("error").Build(),
					},
				},
			},
			nil,
			true,
		},
		{
			"should marshal a conditional step with a default case",
			fields{
				"my_label",
				[]ConditionalCase{
					{
						"this.a == b",
						[]inventory.LogicStep{
							Log(LogLevelInfo).Message("ok").Build(),
						},
					},
				},
				&ConditionalCase{
					"", []inventory.LogicStep{
						Log(LogLevelInfo).Message("error").Build(),
					},
				},
			},
			map[string]any{
				"label": "my_label",
				"switch": []map[string]any{
					{
						"check": "this.a == b",
						"processors": []map[string]any{
							{
								"log": map[string]any{
									"level":   string(LogLevelInfo),
									"message": "ok",
								},
							},
						},
					},
					{
						"processors": []map[string]any{
							{
								"log": map[string]any{
									"level":   string(LogLevelInfo),
									"message": "error",
								},
							},
						},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ConditionalLogicStep{
				label:       tt.fields.label,
				Cases:       tt.fields.Cases,
				DefaultCase: tt.fields.DefaultCase,
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
