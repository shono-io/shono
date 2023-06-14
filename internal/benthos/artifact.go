package benthos

import (
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/rs/xid"
	"github.com/shono-io/shono/core"
	"gopkg.in/yaml.v3"
	"time"
)

func NewArtifact(owner core.Reference, t core.ArtifactType, content map[string]any) (core.Artifact, error) {
	b, err := yaml.Marshal(content)
	if err != nil {
		return nil, err
	}

	// -- read the generated yaml into a streamsBuilder to make sure it is a valid benthos config
	sb := service.NewStreamBuilder()
	if err := sb.SetYAML(string(b)); err != nil {
		return nil, err
	}
	if _, err := sb.Build(); err != nil {
		return nil, err
	}

	return &Artifact{
		Spec: ArtifactSpec{
			Owner:     owner,
			Key:       xid.New().String(),
			Timestamp: time.Now(),
			Type:      t,
			Content:   b,
		},
	}, nil
}

type ArtifactSpec struct {
	Owner     core.Reference
	Key       string
	Timestamp time.Time
	Type      core.ArtifactType
	Content   []byte
}

type Artifact struct {
	Spec ArtifactSpec
}

func (a *Artifact) Owner() core.Reference {
	return a.Spec.Owner
}

func (a *Artifact) Key() string {
	return a.Spec.Key
}

func (a *Artifact) Timestamp() time.Time {
	return a.Spec.Timestamp
}

func (a *Artifact) Type() core.ArtifactType {
	return a.Spec.Type
}

func (a *Artifact) Content() ([]byte, error) {
	return a.Spec.Content, nil
}
