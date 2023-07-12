package dsl

import (
	"github.com/shono-io/shono/inventory"
	"reflect"
	"testing"
)

func TestCatchLogicStep_MarshalBenthos(t *testing.T) {
	type fields struct {
		label string
		Steps []inventory.LogicStep
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{"should marshall a catch step with label",
			fields{
				"my_label",
				[]inventory.LogicStep{
					Log(LogLevelInfo).Message("error").Build(),
				},
			},
			map[string]any{
				"label": "my_label",
				"catch": []map[string]any{
					{
						"log": map[string]any{
							"level":   string(LogLevelInfo),
							"message": "error",
						},
					},
				},
			},
			false,
		},
		{"should marshall a catch step without a label",
			fields{
				"",
				[]inventory.LogicStep{
					Log(LogLevelInfo).Message("error").Build(),
				},
			},
			map[string]any{
				"catch": []map[string]any{
					{
						"log": map[string]any{
							"level":   string(LogLevelInfo),
							"message": "error",
						},
					},
				},
			},
			false,
		},
		{
			"should fail to marshall a catch step without processors",
			fields{
				"",
				[]inventory.LogicStep{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := CatchLogicStep{
				label: tt.fields.label,
				Steps: tt.fields.Steps,
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
