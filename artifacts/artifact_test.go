package artifacts

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewArtifactReference(t *testing.T) {
	type args struct {
		scopeCode    string
		artifactCode string
	}
	tests := []struct {
		name string
		args args
		want commons.Reference
	}{
		{"should create a valid artifact reference", args{"my_scope", "my_artifact"}, "scopes/my_scope/artifacts/my_artifact"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewArtifactReference(tt.args.scopeCode, tt.args.artifactCode); got != tt.want {
				t.Errorf("NewArtifactReference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewArtifactValidation(t *testing.T) {
	gl := GeneratedLogic{
		Steps: []map[string]any{
			{"log": map[string]any{"level": "INFO", "message": ""}},
		},
		Tests: nil,
	}

	concept := inventory.NewConcept("my_scope", "my_concept").Build()

	opts := []Opt{
		WithConcept(concept),
		WithLogic(gl),
	}

	type args struct {
		t        commons.ArtifactType
		ownerRef commons.Reference
		code     string
		opt      []Opt
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"should fail with invalid artifact type", args{"invalid", "my_owner", "my_code", opts}, true},
		{"should fail with invalid owner reference", args{commons.ArtifactTypeConcept, "invalid", "my_code", opts}, true},
		{"should fail with invalid code", args{commons.ArtifactTypeConcept, "my_owner", "", opts}, true},
		{"should fail if concept artifact and no concept is set", args{commons.ArtifactTypeConcept, "my_owner", "my_code", []Opt{WithLogic(gl)}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewArtifact(tt.args.t, tt.args.ownerRef, tt.args.code, tt.args.opt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewArtifact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				assert.True(t, got.Ref.IsValid())
				assert.NotEmpty(t, got.Type)
			}
		})
	}

}

func TestWithConcept(t *testing.T) {
	concept := inventory.NewConcept("my_scope", "my_concept").Build()
	art := &Artifact{}
	WithConcept(concept)(art)
	assert.Equal(t, concept, art.Concept)
}

func TestWithLogic(t *testing.T) {
	gl := GeneratedLogic{
		Steps: []map[string]any{
			{"log": map[string]any{"level": "INFO", "message": ""}},
		},
		Tests: nil,
	}
	art := &Artifact{}
	WithLogic(gl)(art)
	assert.Equal(t, gl, art.Logic)
}

func TestWithInputEvents(t *testing.T) {
	event := inventory.NewEvent("my_scope", "my_concept", "my_event").Build()
	art := &Artifact{}
	WithInputEvents(AsGeneratedEvent(*event))(art)
	assert.Equal(t, []inventory.Event{*event}, art.InputEvents)
}

func TestWithOutputEvents(t *testing.T) {
	event := inventory.NewEvent("my_scope", "my_concept", "my_event").Build()
	art := &Artifact{}
	WithOutputEvents(AsGeneratedEvent(*event))(art)
	assert.Equal(t, []inventory.Event{*event}, art.OutputEvents)
}

func TestWithInput(t *testing.T) {
	inp := GeneratedInput{"my_input", "", map[string]any{}, nil}

	art := &Artifact{}
	WithInput(inp)(art)
	assert.Equal(t, inp, art.Input)
}

func TestWithOutput(t *testing.T) {
	out := GeneratedOutput{"my_output", "", map[string]any{}}

	art := &Artifact{}
	WithOutput(out)(art)
	assert.Equal(t, out, art.Output)
}

func TestWithDlq(t *testing.T) {
	dlq := GeneratedOutput{"my_dlq", "", map[string]any{}}

	art := &Artifact{}
	WithDlq(dlq)(art)
	assert.Equal(t, dlq, art.DLQ)
}
