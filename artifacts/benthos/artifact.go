package benthos

import (
	"github.com/rs/xid"
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
	"time"
)

type Opt func(a *ArtifactSpec)

func WithKey(key string) Opt {
	return func(a *ArtifactSpec) {
		a.Key = key
	}
}

func WithConcept(concept *inventory.Concept) Opt {
	return func(a *ArtifactSpec) {
		a.Concept = concept
	}
}

func NewArtifact(owner commons.Reference, t commons.ArtifactType, logic artifacts.GeneratedLogic, input artifacts.GeneratedInput, output inventory.Output, error inventory.Output, storages []artifacts.Storage, opts ...Opt) (*Artifact, error) {
	res := &Artifact{
		Spec: ArtifactSpec{
			Owner:     owner.String(),
			Key:       xid.New().String(),
			Timestamp: time.Now(),
			Type:      t,
			Logic:     logic,
			Input:     input,
			Output:    output,
			Error:     error,
			Storages:  storages,
			Concept:   nil,
		},
	}

	for _, opt := range opts {
		opt(&res.Spec)
	}

	return res, nil
}

type ArtifactSpec struct {
	Owner     string
	Key       string
	Timestamp time.Time
	Type      commons.ArtifactType

	Input  artifacts.GeneratedInput
	Output inventory.Output
	Error  inventory.Output

	Concept *inventory.Concept `yaml:"concept,omitempty"`

	Storages []artifacts.Storage `yaml:"storages,omitempty"`

	InputEvents  []inventory.Event `yaml:"input_events,omitempty"`
	Logic        artifacts.GeneratedLogic
	OutputEvents []inventory.Event `yaml:"output_events,omitempty"`
}

type Artifact struct {
	Spec ArtifactSpec `yaml:",inline"`
}

func (a *Artifact) Logic() artifacts.GeneratedLogic {
	return a.Spec.Logic
}

func (a *Artifact) Input() artifacts.GeneratedInput {
	return a.Spec.Input
}

func (a *Artifact) Output() inventory.Output {
	return a.Spec.Output
}

func (a *Artifact) Error() inventory.Output {
	return a.Spec.Output
}

func (a *Artifact) Owner() commons.Reference {
	r, _ := commons.ParseString(a.Spec.Owner)
	return r
}

func (a *Artifact) Key() string {
	return a.Spec.Key
}

func (a *Artifact) Timestamp() time.Time {
	return a.Spec.Timestamp
}

func (a *Artifact) Type() commons.ArtifactType {
	return a.Spec.Type
}

func (a *Artifact) Reference() commons.Reference {
	return a.Owner().Child("artifacts", a.Key())
}

func (a *Artifact) Storages() []artifacts.Storage {
	return a.Spec.Storages
}

func (a *Artifact) Concept() *inventory.Concept {
	return a.Spec.Concept
}

func (a *Artifact) InputEvents() []inventory.Event {
	return a.Spec.InputEvents
}

func (a *Artifact) OutputEvents() []inventory.Event {
	return a.Spec.OutputEvents
}
