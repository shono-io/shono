package dsl

import (
	"reflect"
	"testing"
)

func TestLogLogicStep_MarshalBenthos(t *testing.T) {
	type fields struct {
		label   string
		Level   LogLevel
		Message string
		Mapping string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{"should marshall a message log step with label",
			fields{
				"my_label",
				LogLevelInfo,
				"error",
				"",
			}, map[string]any{
				"label": "my_label",
				"log": map[string]any{
					"level":   string(LogLevelInfo),
					"message": "error",
				},
			}, false,
		},
		{"should marshall a mapping log step with label",
			fields{
				"my_label",
				LogLevelInfo,
				"",
				"root := this",
			}, map[string]any{
				"label": "my_label",
				"log": map[string]any{
					"level":          string(LogLevelInfo),
					"fields_mapping": "root := this",
				},
			}, false,
		},
		{"should fail if neither message or mapping is provided",
			fields{
				"my_label",
				LogLevelInfo,
				"",
				"",
			},
			nil,
			true,
		},
		{"should fail if an invalid log level is provided",
			fields{
				"my_label",
				"mogus",
				"error",
				"",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := LogLogicStep{
				label:   tt.fields.label,
				Level:   tt.fields.Level,
				Message: tt.fields.Message,
				Mapping: tt.fields.Mapping,
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
