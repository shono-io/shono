package commons

type ArtifactType string

var (
	ArtifactTypeInjector  ArtifactType = "injector"
	ArtifactTypeExtractor ArtifactType = "extractor"
	ArtifactTypeConcept   ArtifactType = "reaktor"
)

// Status of a component.
type Status string

// Node statuses.
var (
	StatusStable       Status = "stable"
	StatusBeta         Status = "beta"
	StatusExperimental Status = "experimental"
	StatusDeprecated   Status = "deprecated"
)
