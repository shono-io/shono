package dsl

import (
	"github.com/shono-io/shono/inventory"
	"reflect"
	"testing"
)

func TestBranchStep_MarshalBenthos(t *testing.T) {
	type fields struct {
		Label string
		Pre   string
		Steps []inventory.LogicStep
		Post  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{
			"should marshall a complete branch step with label",
			fields{
				"my_label",
				`root = this`,
				[]inventory.LogicStep{
					Log(LogLevelInfo).Message("error").Build(),
				},
				`root.result = this`,
			}, map[string]any{
				"label": "my_label",
				"branch": map[string]any{
					"request_map": `root = this`,
					"processors": []map[string]any{
						{
							"log": map[string]any{
								"level":   string(LogLevelInfo),
								"message": "error",
							},
						},
					},
					"response_map": `root.result = this`,
				},
			}, false,
		},
		{
			"should marshall a branch step with only processors",
			fields{
				"",
				"",
				[]inventory.LogicStep{
					Log(LogLevelInfo).Message("error").Build(),
				},
				"",
			}, map[string]any{
				"branch": map[string]any{
					"processors": []map[string]any{
						{
							"log": map[string]any{
								"level":   string(LogLevelInfo),
								"message": "error",
							},
						},
					},
				},
			}, false,
		},
		{
			"should fail if no processors are available",
			fields{
				"my_label",
				`root = this`,
				[]inventory.LogicStep{},
				`root.result = this`,
			}, nil, true,
		},
		{
			"should fail if no processors is nil",
			fields{
				"my_label",
				`root = this`,
				nil,
				`root.result = this`,
			}, nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BranchStep{
				label: tt.fields.Label,
				Pre:   tt.fields.Pre,
				Steps: tt.fields.Steps,
				Post:  tt.fields.Post,
			}
			got, err := b.MarshalBenthos("")
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
