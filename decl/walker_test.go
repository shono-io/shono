package decl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalk(t *testing.T) {
	t.Run("should walk all specs", shouldTriggerHandlerForSpecs)
}

func shouldTriggerHandlerForSpecs(t *testing.T) {
	h := &memHandler{}

	err := Walk("test", h)
	assert.NoError(t, err)

	assert.NotNil(t, h.scope)
	assert.NoError(t, h.scope.Validate())

	assert.NotNil(t, h.concept)
	assert.NoError(t, h.concept.Validate())

	assert.NotNil(t, h.event)
	assert.NoError(t, h.event.Validate())

	assert.NotNil(t, h.reactor)
}

type memHandler struct {
	scope   *ScopeSpec
	concept *ConceptSpec
	event   *EventSpec

	reactor *ReactorSpec
}

func (m *memHandler) OnScope(scope *ScopeSpec) error {
	m.scope = scope
	return nil
}

func (m *memHandler) OnConcept(concept *ConceptSpec) error {
	m.concept = concept
	return nil
}

func (m *memHandler) OnEvent(event *EventSpec) error {
	m.event = event
	return nil
}

func (m *memHandler) OnReactor(reactor *ReactorSpec) error {
	m.reactor = reactor
	return nil
}
