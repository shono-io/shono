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

func (g *ConceptGenerator) Generate(applicationId string, artifactId string, inv inventory.Inventory, conceptRef commons.Reference) (*artifacts.Artifact, error) {
	concept, err := inv.ResolveConcept(conceptRef)
	if err != nil {
		return nil, err
	}
	opts := []artifacts.Opt{
		artifacts.WithConcept(concept),
	}

	reactors, err := inv.ListReactorsForConcept(conceptRef)
	if err != nil {
		return nil, err
	}

	var inputEventRefs []commons.Reference
	var inputEvents []artifacts.GeneratedEvent
	var outputEvents []artifacts.GeneratedEvent
	for _, r := range reactors {
		inputEventRefs = append(inputEventRefs, r.InputEvent)
		evt, err := inv.ResolveEvent(r.InputEvent)
		if err != nil {
			return nil, fmt.Errorf("input event %q: %w", r.InputEvent, err)
		}
		if evt == nil {
			return nil, fmt.Errorf("input event: %s not found", r.InputEvent)
		}
		inputEvents = append(inputEvents, artifacts.AsGeneratedEvent(*evt))

		uniqueOutputEvents := make(map[commons.Reference]struct{})
		for _, outEvt := range r.OutputEventCodes {
			uniqueOutputEvents[r.Concept.Child("events", outEvt)] = struct{}{}
		}
		for oe, _ := range uniqueOutputEvents {
			evt, err := inv.ResolveEvent(oe)
			if err != nil {
				return nil, fmt.Errorf("output event: %w", err)
			}
			if evt == nil {
				return nil, fmt.Errorf("output event: %s not found", outputEvents)
			}
			outputEvents = append(outputEvents, artifacts.AsGeneratedEvent(*evt))
		}
	}
	opts = append(opts,
		artifacts.WithInputEvents(inputEvents...),
		artifacts.WithOutputEvents(outputEvents...))

	inp, err := generateBackboneInput(applicationId, inputEventRefs)
	if err != nil {
		return nil, fmt.Errorf("input: %w", err)
	}
	opts = append(opts, artifacts.WithInput(*inp))

	logic, err := generateWrapperLogic(reactors)
	if err != nil {
		return nil, fmt.Errorf("logic: %w", err)
	}
	l, err := generateLogic(logic)
	if err != nil {
		return nil, fmt.Errorf("logic: %w", err)
	}
	opts = append(opts, artifacts.WithLogic(*l))

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

	return artifacts.NewArtifact(commons.ArtifactTypeConcept, conceptRef, fmt.Sprintf("%s_%s", applicationId, artifactId), opts...)
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
