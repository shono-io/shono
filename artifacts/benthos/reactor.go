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

func (g *ConceptGenerator) Generate(applicationId string, artifactId string, inv inventory.Inventory, conceptRef commons.Reference) (artifacts.Artifact, error) {
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
		inputEvents = append(inputEvents, r.InputEvent)
	}

	inp, err := generateBackboneInput(applicationId, inputEvents)
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

	l, err := generateLogic(logic)
	if err != nil {
		return nil, fmt.Errorf("logic: %w", err)
	}

	return NewArtifact(conceptRef, commons.ArtifactTypeConcept, *l, *inp, out, dlq, WithConcept(concept), WithKey(artifactId))
}

func generateWrapperLogic(reactors []inventory.Reactor) (inventory.Logic, error) {
	result := inventory.NewLogic()
	var cases []dsl.ConditionalCase
	for _, r := range reactors {
		cases = append(cases, dsl.SwitchCase(
			fmt.Sprintf("@shono_kind == %q", r.InputEvent.String()),
			r.Logic.Steps()...))

		// -- add the tests
		result.Test(r.Logic.Tests()...)
	}

	// -- add a default case that logs unmatched event to trace
	cases = append(cases, dsl.SwitchDefault(
		dsl.Log("TRACE", "no processor for ${!meta(\"kind\")} with payload ${!this.format_json()}"),
		dsl.Transform(dsl.BloblangMapping(`root = deleted()`)),
	))

	return result.
		Steps(
			dsl.Log("TRACE", "received ${!@} with payload ${!this.format_json()}"),
			dsl.Switch(cases...),
		).
		Build(), nil
}
