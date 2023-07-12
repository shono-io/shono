package decl

import (
	"fmt"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/local"
)

func NewInventory(paths ...string) (inventory.Inventory, error) {
	h := &handler{
		inv: local.NewInventory(),
	}

	for _, path := range paths {
		if err := Walk(path, h); err != nil {
			return nil, err
		}
	}

	return h.inv.Build()
}

type handler struct {
	inv *local.InventoryBuilder
}

func (h *handler) OnScope(scope *ScopeSpec) error {
	status, err := commons.StatusOf(scope.Status)
	if err != nil {
		return err
	}

	res := inventory.NewScope(scope.Code).
		Summary(scope.Summary).
		Status(status).
		Docs(scope.Docs).
		Build()

	h.inv.Scope(res)

	return nil
}

func (h *handler) OnConcept(concept *ConceptSpec) error {
	status, err := commons.StatusOf(concept.Status)
	if err != nil {
		return err
	}

	res := inventory.NewConcept(concept.Scope, concept.Code).
		Summary(concept.Summary).
		Status(status).
		Docs(concept.Docs).
		IsStored(concept.Stored).
		Build()

	h.inv.Concept(res)

	return nil
}

func (h *handler) OnEvent(event *EventSpec) error {
	status, err := commons.StatusOf(event.Status)
	if err != nil {
		return err
	}

	res := inventory.NewEvent(event.Scope, event.Concept, event.Code).
		Summary(event.Summary).
		Status(status).
		Docs(event.Docs).
		Build()

	h.inv.Event(res)

	return nil
}

func (h *handler) OnReactor(reactor *ReactorSpec) error {
	if reactor.Status == "" {
		reactor.Status = "experimental"
	}

	status, err := commons.StatusOf(reactor.Status)
	if err != nil {
		return err
	}

	conceptRef := inventory.NewConceptReference(reactor.Concept.Scope, reactor.Concept.Code)

	var steps []inventory.StepBuilder
	for idx, step := range reactor.Logic {
		st := step.AsLogicStep(fmt.Sprintf("step[%d]", idx), conceptRef)
		steps = append(steps, st)
	}

	lb := inventory.NewLogic().Steps(steps...)

	rb := inventory.NewReactor(reactor.Concept.Scope, reactor.Concept.Code, fmt.Sprintf("on_%s_%s_%s", reactor.InputEvent.Scope, reactor.InputEvent.Concept, reactor.InputEvent.Code)).
		Summary(reactor.Summary).
		Status(status).
		Docs(reactor.Docs).
		InputEvent(inventory.NewEventReference(reactor.InputEvent.Scope, reactor.InputEvent.Concept, reactor.InputEvent.Code)).
		OutputEventCodes(reactor.OutputEventCodes()...).
		Logic(lb)

	h.inv.Reactor(rb.Build())

	return nil
}
