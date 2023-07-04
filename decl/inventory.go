package decl

import (
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
