package graph

// Status of a component.
type Status string

// Node statuses.
var (
	StatusStable       Status = "stable"
	StatusBeta         Status = "beta"
	StatusExperimental Status = "experimental"
	StatusDeprecated   Status = "deprecated"
)

// Type of node.
type Type string

var (
	TypeScope   Type = "scope"
	TypeConcept Type = "concept"
	TypeEvent   Type = "event"
	TypeReaktor Type = "reaktor"
)

var Types = []Type{
	TypeScope,
	TypeConcept,
	TypeEvent,
	TypeReaktor,
}

// NodeSpec describes a Shono node. The idea has been copied from Benthos' ComponentSpec.
type NodeSpec struct {
	// Code is the unique identifier of the node.
	Code string `json:"code"`

	// Type is the type of the node.
	Type Type `json:"type"`

	// Status is the status of the node.
	Status Status `json:"status"`

	// Summary is a short summary of the node.
	Summary string `json:"summary,omitempty"`

	// Docs is the documentation for the node.
	Docs string `json:"docs,omitempty"`

	// Version is the version of shono in which the node was introduced.
	Version string `json:"version,omitempty"`
}

func NewNodeSpec(code string, t Type) *NodeSpec {
	return &NodeSpec{
		Code:   code,
		Type:   t,
		Status: StatusExperimental,
	}
}

func (s *NodeSpec) WithSummary(summary string) *NodeSpec {
	s.Summary = summary
	return s
}

func (s *NodeSpec) WithDocs(docs string) *NodeSpec {
	s.Docs = docs
	return s
}

func (s *NodeSpec) Since(version string) *NodeSpec {
	s.Version = version
	return s
}

func (s *NodeSpec) Stable() *NodeSpec {
	s.Status = StatusStable
	return s
}

func (s *NodeSpec) Beta() *NodeSpec {
	s.Status = StatusBeta
	return s
}

func (s *NodeSpec) Experimental() *NodeSpec {
	s.Status = StatusExperimental
	return s
}
