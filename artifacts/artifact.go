package artifacts

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
	"time"
)

func NewArtifactReference(scopeCode, artifactCode string) commons.Reference {
	return inventory.NewScopeReference(scopeCode).Child("artifacts", artifactCode)
}

type GeneratedLogic struct {
	Steps []map[string]any `yaml:"steps"`
	Tests []map[string]any `yaml:"tests,omitempty"`
}

type Artifact interface {
	Owner() commons.Reference
	Key() string
	Timestamp() time.Time
	Type() commons.ArtifactType
	Reference() commons.Reference

	Logic() GeneratedLogic

	Input() inventory.Input
	Output() inventory.Output
	Error() inventory.Output

	Storages() []inventory.Storage
}

type ArtifactLoader interface {
	LoadArtifact(ref commons.Reference) (Artifact, error)
}

type ArtifactDumper interface {
	StoreArtifact(artifact Artifact) error
}
