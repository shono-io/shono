package local

import (
	"github.com/shono-io/shono/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConceptRepo_ListConceptsForScope(t *testing.T) {
	registry := NewRegistry(
		WithScope(graph.NewScope("my-scope")),
		WithConcept(*graph.NewConcept("my-scope", "my-concept")),
	)

	l, err := registry.ListConceptsForScope("my-scope")
	assert.NoError(t, err)
	assert.Equal(t, []graph.Concept{*graph.NewConcept("my-scope", "my-concept")}, l)
}
