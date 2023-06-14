package core

import "time"

type ArtifactType string

var (
	ArtifactTypeInjector  ArtifactType = "injector"
	ArtifactTypeExtractor ArtifactType = "extractor"
	ArtifactTypeReaktor   ArtifactType = "reaktor"
)

type Artifact interface {
	Owner() Reference
	Key() string
	Timestamp() time.Time
	Type() ArtifactType
	Content() ([]byte, error)
}
