package benthos

import (
	"fmt"
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
)

func NewInjectorGenerator() *InjectorGenerator {
	return &InjectorGenerator{}
}

type InjectorGenerator struct {
}

func (i *InjectorGenerator) Generate(applicationId string, artifactId string, inv inventory.Inventory, injectorRef commons.Reference) (artifacts.Artifact, error) {
	injector, err := inv.ResolveInjector(injectorRef)
	if err != nil {
		return nil, err
	}

	var inpLogic *artifacts.GeneratedLogic
	if injector.Input.Logic != nil {
		l, err := generateLogic(injector.Input.Logic)
		if err != nil {
			return nil, err
		}
		inpLogic = l
	}
	inp := artifacts.GeneratedInput{
		Id:         injector.Input.Id,
		Kind:       injector.Input.Kind,
		ConfigSpec: injector.Input.ConfigSpec,
		Config:     injector.Input.Config,
		Logic:      inpLogic,
	}

	out, err := generateBackboneOutput()
	if err != nil {
		return nil, fmt.Errorf("output: %w", err)
	}

	dlq, err := generateBackboneDLQ()
	if err != nil {
		return nil, fmt.Errorf("dlq: %w", err)
	}

	logic, err := generateLogic(injector.Logic)
	if err != nil {
		return nil, fmt.Errorf("logic: %w", err)
	}

	return NewArtifact(injectorRef.Parent(), commons.ArtifactTypeInjector, *logic, inp, out, dlq, WithKey(artifactId))
}
