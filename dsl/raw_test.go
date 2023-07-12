package dsl

import (
	"reflect"
	"testing"
)

func TestRawLogicStep_MarshalBenthos(t *testing.T) {
	type fields struct {
		label   string
		Content map[string]any
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]any
		wantErr bool
	}{
		{
			name: "should marshal with label",
			fields: fields{
				label: "test",
				Content: map[string]any{
					"http": map[string]any{
						"url": "http://localhost:8080",
					},
				},
			},
			want: map[string]any{
				"label": "test",
				"http": map[string]any{
					"url": "http://localhost:8080",
				},
			},
		},
		{
			name: "should marshal without label",
			fields: fields{
				Content: map[string]any{
					"http": map[string]any{
						"url": "http://localhost:8080",
					},
				},
			},
			want: map[string]any{
				"http": map[string]any{
					"url": "http://localhost:8080",
				},
			},
		},
		{
			name: "should fail if content is nil",
			fields: fields{
				Content: nil,
			},
			wantErr: true,
		},
		{
			name: "should fail if content is empty",
			fields: fields{
				Content: map[string]any{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := RawLogicStep{
				label:   tt.fields.label,
				Content: tt.fields.Content,
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
