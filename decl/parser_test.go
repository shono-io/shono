package decl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Testparser(t *testing.T) {
	t.Run("should parse full scope", shouldParseFullScope)
	t.Run("should parse partial scope", shouldParsePartialScope)
	t.Run("should fail to parse scope with missing code", shouldFailToParseScopeWithMissingCode)
	t.Run("should fail to parse scope with missing summary", shouldFailToParseScopeWithMissingSummary)

	t.Run("should parse full concept", shouldParseFullConcept)
	t.Run("should parse partial concept", shouldParsePartialConcept)
	t.Run("should fail to parse concept with missing scope", shouldFailToParseConceptWithMissingScope)
	t.Run("should fail to parse concept with missing code", shouldFailToParseConceptWithMissingCode)
	t.Run("should fail to parse concept with missing summary", shouldFailToParseConceptWithMissingSummary)

	t.Run("should parse full event", shouldParseFullEvent)
	t.Run("should parse partial event", shouldParsePartialEvent)
	t.Run("should fail to parse event with missing scope", shouldFailToParseEventWithMissingScope)
	t.Run("should fail to parse event with missing concept", shouldFailToParseEventWithMissingConcept)
	t.Run("should fail to parse event with missing code", shouldFailToParseEventWithMissingCode)
	t.Run("should fail to parse event with missing summary", shouldFailToParseEventWithMissingSummary)
}

func shouldParseFullScope(t *testing.T) {
	content := []byte(`
scope:
  code: my_scope
  status: experimental
  summary: my_summary
  docs: A simple scope
`)

	res, err := (&parser{}).Parse(content)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Scope)
	assert.Equal(t, "my_scope", res.Scope.Code)
	assert.Equal(t, "experimental", res.Scope.Status)
	assert.Equal(t, "my_summary", res.Scope.Summary)
	assert.Equal(t, "A simple scope", res.Scope.Docs)
}

func shouldParsePartialScope(t *testing.T) {
	content := []byte(`
scope:
  code: my_scope
  summary: my_summary
`)

	res, err := (&parser{}).Parse(content)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Scope)
	assert.Equal(t, "my_scope", res.Scope.Code)
	assert.Equal(t, "experimental", res.Scope.Status)
	assert.Equal(t, "my_summary", res.Scope.Summary)
	assert.Equal(t, "", res.Scope.Docs)
}

func shouldFailToParseScopeWithMissingCode(t *testing.T) {
	content := []byte(`
scope:
  summary: my_summary
`)
	res, err := (&parser{}).Parse(content)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func shouldFailToParseScopeWithMissingSummary(t *testing.T) {
	content := []byte(`
scope:
  code: my_scope
`)
	res, err := (&parser{}).Parse(content)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func shouldParseFullConcept(t *testing.T) {
	content := []byte(`
concept:
  scope: my_scope
  code: my_concept
  summary: my_summary
  status: experimental
  docs: A simple concept
  stored: true
`)
	res, err := (&parser{}).Parse(content)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Concept)
	assert.Equal(t, "my_scope", res.Concept.Scope)
	assert.Equal(t, "my_concept", res.Concept.Code)
	assert.Equal(t, "experimental", res.Concept.Status)
	assert.Equal(t, "my_summary", res.Concept.Summary)
	assert.Equal(t, "A simple concept", res.Concept.Docs)
	assert.Equal(t, true, res.Concept.Stored)
}

func shouldParsePartialConcept(t *testing.T) {
	content := []byte(`
concept:
  scope: my_scope
  code: my_concept
  summary: my_summary
`)
	res, err := (&parser{}).Parse(content)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Concept)
	assert.Equal(t, "my_scope", res.Concept.Scope)
	assert.Equal(t, "my_concept", res.Concept.Code)
	assert.Equal(t, "experimental", res.Concept.Status)
	assert.Equal(t, "my_summary", res.Concept.Summary)
	assert.Equal(t, "", res.Concept.Docs)
	assert.Equal(t, false, res.Concept.Stored)
}

func shouldFailToParseConceptWithMissingScope(t *testing.T) {
	content := []byte(`
concept:
  code: my_concept
  summary: my_summary
`)
	res, err := (&parser{}).Parse(content)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func shouldFailToParseConceptWithMissingCode(t *testing.T) {
	content := []byte(`
concept:
  scope: my_scope
  summary: my_summary
`)
	res, err := (&parser{}).Parse(content)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func shouldFailToParseConceptWithMissingSummary(t *testing.T) {
	content := []byte(`
concept:
  scope: my_scope
  code: my_event
`)
	res, err := (&parser{}).Parse(content)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func shouldParseFullEvent(t *testing.T) {
	content := []byte(`
event:
  scope: my_scope
  concept: my_concept
  code: my_event
  status: experimental
  summary: my_summary
  docs: A simple event
`)

	res, err := (&parser{}).Parse(content)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Event)
	assert.Equal(t, "my_scope", res.Event.Scope)
	assert.Equal(t, "my_concept", res.Event.Concept)
	assert.Equal(t, "my_event", res.Event.Code)
	assert.Equal(t, "experimental", res.Event.Status)
	assert.Equal(t, "my_summary", res.Event.Summary)
	assert.Equal(t, "A simple event", res.Event.Docs)
}

func shouldParsePartialEvent(t *testing.T) {
	content := []byte(`
event:
  scope: my_scope
  concept: my_concept
  code: my_event
  summary: my_summary
`)

	res, err := (&parser{}).Parse(content)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Event)
	assert.Equal(t, "my_scope", res.Event.Scope)
	assert.Equal(t, "my_concept", res.Event.Concept)
	assert.Equal(t, "my_event", res.Event.Code)
	assert.Equal(t, "experimental", res.Event.Status)
	assert.Equal(t, "my_summary", res.Event.Summary)
	assert.Equal(t, "", res.Event.Docs)
}

func shouldFailToParseEventWithMissingScope(t *testing.T) {
	content := []byte(`
event:
  concept: my_concept	
  code: my_event	
  summary: my_summary
`)
	res, err := (&parser{}).Parse(content)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func shouldFailToParseEventWithMissingConcept(t *testing.T) {
	content := []byte(`
event:
  scope: my_scope
  code: my_event
  summary: my_summary
`)
	res, err := (&parser{}).Parse(content)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func shouldFailToParseEventWithMissingCode(t *testing.T) {
	content := []byte(`
event:
  scope: my_scope
  concept: my_concept
  summary: my_summary
`)
	res, err := (&parser{}).Parse(content)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func shouldFailToParseEventWithMissingSummary(t *testing.T) {
	content := []byte(`
event:
  scope: my_scope
  concept: my_concept
  code: my_event
`)
	res, err := (&parser{}).Parse(content)
	assert.Error(t, err)
	assert.Nil(t, res)
}
