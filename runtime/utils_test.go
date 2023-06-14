package runtime

import (
	"github.com/shono-io/shono/core"
	"testing"
)

func Test_labelize(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should convert simple keys", args{core.NewReference("org", "one").String()}, "org_one"},
		{"should convert simple keys", args{core.NewReference("org", "one").Child("scope", "two").String()}, "org_one_scope_two"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := labelize(tt.args.s); got != tt.want {
				t.Errorf("labelize() = %v, want %v", got, tt.want)
			}
		})
	}
}
