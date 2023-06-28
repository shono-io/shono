package artifacts

import (
	"fmt"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
	"time"
)

func NewArtifactReference(scopeCode, artifactCode string) commons.Reference {
	return inventory.NewScopeReference(scopeCode).Child("artifacts", artifactCode)
}

type Opt func(a *Artifact)

func WithConcept(concept *inventory.Concept) Opt {
	return func(a *Artifact) {
		a.Concept = concept
	}
}

func WithLogic(logic GeneratedLogic) Opt {
	return func(a *Artifact) {
		a.Logic = logic
	}
}

func WithInputEvents(events ...GeneratedEvent) Opt {
	return func(a *Artifact) {
		a.InputEvents = append(a.InputEvents, events...)
	}
}

func WithOutputEvents(events ...GeneratedEvent) Opt {
	return func(a *Artifact) {
		a.OutputEvents = append(a.OutputEvents, events...)
	}
}

func WithInput(input GeneratedInput) Opt {
	return func(a *Artifact) {
		a.Input = input
	}
}

func WithOutput(output GeneratedOutput) Opt {
	return func(a *Artifact) {
		a.Output = output
	}
}

func WithDlq(dlq GeneratedOutput) Opt {
	return func(a *Artifact) {
		a.DLQ = dlq
	}
}

func NewArtifact(t commons.ArtifactType, ownerRef commons.Reference, code string, opt ...Opt) (*Artifact, error) {
	if t != commons.ArtifactTypeConcept && t != commons.ArtifactTypeInjector && t != commons.ArtifactTypeExtractor {
		return nil, fmt.Errorf("invalid artifact type")
	}

	if code == "" {
		return nil, fmt.Errorf("artifact code cannot be empty")
	}

	if !ownerRef.IsValid() {
		return nil, fmt.Errorf("artifact owner reference is invalid")
	}

	art := &Artifact{
		Ref:       ownerRef.Child("artifacts", code),
		Type:      t,
		Timestamp: time.Now(),
	}

	for _, o := range opt {
		o(art)
	}

	// -- perform validations
	switch art.Type {
	case commons.ArtifactTypeConcept:
		if art.Concept == nil {
			return nil, fmt.Errorf("concept artifact must have a concept")
		}
	}

	if art.Logic.Steps == nil {
		return nil, fmt.Errorf("logic must have at least one step")
	}

	return art, nil
}

type Artifact struct {
	Type      commons.ArtifactType `yaml:"type"`
	Ref       commons.Reference    `yaml:"ref"`
	Timestamp time.Time            `yaml:"timestamp"`

	Concept *inventory.Concept `yaml:"concept,omitempty"`

	Logic GeneratedLogic `yaml:"logic,omitempty"`

	InputEvents  []GeneratedEvent `yaml:"input_events,omitempty"`
	OutputEvents []GeneratedEvent `yaml:"output_events,omitempty"`

	Input  GeneratedInput  `yaml:"input"`
	Output GeneratedOutput `yaml:"output"`
	DLQ    GeneratedOutput `yaml:"error,omitempty"`
}

type GeneratedLogic struct {
	Steps []map[string]any `yaml:"steps"`
	Tests []map[string]any `yaml:"tests,omitempty"`
}

func AsGeneratedEvent(evt inventory.Event) GeneratedEvent {
	return GeneratedEvent{
		Ref:     evt.Reference(),
		Summary: evt.Summary,
		Docs:    evt.Docs,
		Status:  evt.Status,
	}
}

type GeneratedEvent struct {
	Ref     commons.Reference `yaml:"ref"`
	Summary string            `yaml:"summary"`
	Docs    string            `yaml:"docs"`
	Status  commons.Status    `yaml:"status"`
}

type ArtifactLoader interface {
	LoadArtifact(ref commons.Reference) (Artifact, error)
}

type ArtifactDumper interface {
	StoreArtifact(artifact Artifact) error
}
