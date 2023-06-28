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

func (i *InjectorGenerator) Generate(applicationId string, artifactId string, inv inventory.Inventory, injectorRef commons.Reference) (*artifacts.Artifact, error) {
	injector, err := inv.ResolveInjector(injectorRef)
	if err != nil {
		return nil, err
	}

	var outputEvents []artifacts.GeneratedEvent
	for _, e := range injector.OutputEvents {
		evt, err := inv.ResolveEvent(e)
		if err != nil {
			return nil, fmt.Errorf("output event: %w", err)
		}
		if evt == nil {
			return nil, fmt.Errorf("output event: %s not found", e)
		}
		outputEvents = append(outputEvents, artifacts.AsGeneratedEvent(*evt))
	}
	opts := []artifacts.Opt{
		artifacts.WithOutputEvents(outputEvents...),
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
		Id:     injector.Input.Id,
		Kind:   injector.Input.Kind,
		Config: injector.Input.Config,
		Logic:  inpLogic,
	}
	opts = append(opts, artifacts.WithInput(inp))

	out, err := generateBackboneOutput()
	if err != nil {
		return nil, fmt.Errorf("output: %w", err)
	}
	opts = append(opts, artifacts.WithOutput(*out))

	dlq, err := generateBackboneDLQ()
	if err != nil {
		return nil, fmt.Errorf("dlq: %w", err)
	}
	opts = append(opts, artifacts.WithDlq(*dlq))

	logic, err := generateLogic(injector.Logic)
	if err != nil {
		return nil, fmt.Errorf("logic: %w", err)
	}
	opts = append(opts, artifacts.WithLogic(*logic))

	return artifacts.NewArtifact(commons.ArtifactTypeInjector, injector.Reference().Parent(), fmt.Sprintf("%s_%s", applicationId, artifactId), opts...)
}
