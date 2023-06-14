package core

// Status of a component.
type Status string

// Node statuses.
var (
	StatusStable       Status = "stable"
	StatusBeta         Status = "beta"
	StatusExperimental Status = "experimental"
	StatusDeprecated   Status = "deprecated"
)

type Node interface {
	Code() string
	Summary() string
	Docs() string
	Status() Status
	Reference() Reference
}

type Executable interface {
	Node
	Logic() Logic
}
