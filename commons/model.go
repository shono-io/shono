package commons

import "fmt"

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

func StatusOf(s string) (Status, error) {
	switch s {
	case "stable":
		return StatusStable, nil
	case "beta":
		return StatusBeta, nil
	case "experimental":
		return StatusExperimental, nil
	case "deprecated":
		return StatusDeprecated, nil
	default:
		return "", fmt.Errorf("invalid status: %s", s)
	}
}
