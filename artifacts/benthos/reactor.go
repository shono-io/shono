package benthos

import (
	"fmt"
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/inventory"
)

func NewConceptGenerator() *ConceptGenerator {
	return &ConceptGenerator{}
}

type ConceptGenerator struct {
}

func (g *ConceptGenerator) Generate(artifactId string, inv inventory.Inventory, conceptRef commons.Reference) (artifacts.Artifact, error) {
	concept, err := inv.ResolveConcept(conceptRef)
	if err != nil {
		return nil, err
	}

	reactors, err := inv.ListReactorsForConcept(conceptRef)
	if err != nil {
		return nil, err
	}

	var inputEvents []commons.Reference
	for _, r := range reactors {
		inputEvents = append(inputEvents, r.InputEvent())
	}

	inp, err := generateBackboneInput(inputEvents)
	if err != nil {
		return nil, fmt.Errorf("input: %w", err)
	}

	logic, err := generateWrapperLogic(reactors)
	if err != nil {
		return nil, fmt.Errorf("logic: %w", err)
	}

	out, err := generateBackboneOutput()
	if err != nil {
		return nil, fmt.Errorf("output: %w", err)
	}

	dlq, err := generateBackboneDLQ()
	if err != nil {
		return nil, fmt.Errorf("dlq: %w", err)
	}

	var storages []artifacts.Storage
	if concept.Stored() {
		storages = append(storages, artifacts.Storage{Collection: fmt.Sprintf("%s__%s", conceptRef.Parent().Code(), conceptRef.Code())})
	}

	l, err := generateLogic(logic)
	if err != nil {
		return nil, fmt.Errorf("logic: %w", err)
	}

	return NewArtifact(conceptRef, commons.ArtifactTypeConcept, *l, inp, out, dlq, storages, WithConcept(&concept), WithKey(artifactId))
}

func generateWrapperLogic(reactors []inventory.Reactor) (inventory.Logic, error) {
	result := inventory.NewLogic()
	var cases []dsl.ConditionalCase
	for _, r := range reactors {
		cases = append(cases, dsl.SwitchCase(
			fmt.Sprintf("@io_shono_kind == %q", r.InputEvent().String()),
			r.Logic().Steps()...))

		// -- add the tests
		result.Test(r.Logic().Tests()...)
	}

	// -- add a default case that logs unmatched event to trace
	cases = append(cases, dsl.SwitchDefault(dsl.Log("TRACE", "no processor for ${!this.format_json()}")))

	return result.
		Steps(
			dsl.Switch(cases...),
		).
		Build(), nil
}
