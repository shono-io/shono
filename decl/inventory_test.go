package decl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInventory(t *testing.T) {
	t.Run("should create new inventory based on directory contents", shouldCreateNewInventoryBasedOnDirectoryContents)
}

func shouldCreateNewInventoryBasedOnDirectoryContents(t *testing.T) {
	inv, err := NewInventory("test")
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, inv)
	scopes, err := inv.ListScopes()
	assert.NoError(t, err)
	assert.Len(t, scopes, 1)

	concepts, err := inv.ListConceptsForScope(scopes[0].Reference())
	assert.NoError(t, err)
	assert.Len(t, concepts, 1)

	events, err := inv.ListEventsForConcept(concepts[0].Reference())
	assert.NoError(t, err)
	assert.Len(t, events, 1)
}
