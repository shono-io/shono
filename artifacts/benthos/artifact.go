package benthos

import (
	"github.com/rs/xid"
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
	"time"
)

func NewArtifact(owner commons.Reference, t commons.ArtifactType, concept *inventory.Concept, logic artifacts.GeneratedLogic, input inventory.Input, output inventory.Output, error inventory.Output, storages []artifacts.Storage) (*Artifact, error) {
	return &Artifact{
		Spec: ArtifactSpec{
			Owner:     owner,
			Key:       xid.New().String(),
			Timestamp: time.Now(),
			Type:      t,
			Logic:     logic,
			Input:     input,
			Output:    output,
			Error:     error,
			Storages:  storages,
			Concept:   concept,
		},
	}, nil
}

type ArtifactSpec struct {
	Owner     commons.Reference
	Key       string
	Timestamp time.Time
	Type      commons.ArtifactType

	Input  inventory.Input
	Output inventory.Output
	Error  inventory.Output

	Concept *inventory.Concept

	Storages []artifacts.Storage

	Logic artifacts.GeneratedLogic
}

type Artifact struct {
	Spec ArtifactSpec
}

func (a *Artifact) Logic() artifacts.GeneratedLogic {
	return a.Spec.Logic
}

func (a *Artifact) Input() inventory.Input {
	return a.Spec.Input
}

func (a *Artifact) Output() inventory.Output {
	return a.Spec.Output
}

func (a *Artifact) Error() inventory.Output {
	return a.Spec.Output
}

func (a *Artifact) Owner() commons.Reference {
	return a.Spec.Owner
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
