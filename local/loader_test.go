package local

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadBytes(t *testing.T) {
	var content = `
backbone:
  kind: kafka
  config:
    bootstrap_servers:
      - pkc-1wvvj.westeurope.azure.confluent.cloud:9092
    sasl:
      - mechanism: PLAIN
        username: ${KAFKA_API_KEY}
        password: ${KAFKA_API_SECRET}

storages:
  my_mongo:
    kind: mongodb
    config:
      uri: ${MDB_URI}
      database: "shono-dev"

scopes:
  todo:
    name: Todo
    description: Unit responsible for managing todo lists
    concepts:
      task:
        name: Task
        description: A task to be completed
        single: task
        plural: tasks
        events:
          operation_failed:
            name: Operation Failed
            description: An operation on a task failed
`

	os.Setenv("KAFKA_API_KEY", "my-key")
	os.Setenv("KAFKA_API_SECRET", "my-secret")
	os.Setenv("MDB_URI", "mongodb://localhost:27017")

	reg, err := LoadBytes([]byte(content))
	if err != nil {
		t.Fatal(err)
	}

	// -- we expect one backbone
	bb, err := reg.GetBackbone()
	assert.NoError(t, err)
	assert.NotNil(t, bb)

	// -- we expect one storage
	storages := reg.ListStorages()
	assert.Len(t, storages, 1)

	storage := reg.GetStorage("my_mongo")
	assert.NotNil(t, storage)
	assert.Equal(t, "my_mongo", storage.Key())

	// -- we expect one scope
	scopes, err := reg.ListScopes()
	assert.NoError(t, err)
	assert.Len(t, scopes, 1)

	scope, err := reg.GetScope("todo")
	assert.NoError(t, err)
	assert.NotNil(t, scope)
	assert.Equal(t, "todo", scope.Code)
	assert.Equal(t, "Todo", scope.Name)
	assert.Equal(t, "Unit responsible for managing todo lists", scope.Description)

	// -- we expect one concept within the scope
	concepts, err := reg.ListConceptsForScope(scope.Code)
	assert.NoError(t, err)
	assert.Len(t, concepts, 1)

	concept, err := reg.GetConcept(scope.Code, "task")
	assert.NoError(t, err)
	assert.NotNil(t, concept)
	assert.Equal(t, "task", concept.Code)
	assert.Equal(t, "Task", concept.Name)
	assert.Equal(t, "A task to be completed", concept.Description)
	assert.Equal(t, "task", concept.Single)
	assert.Equal(t, "tasks", concept.Plural)

	// -- we expect one event within the concept
	events, err := reg.ListEventsForConcept(scope.Code, concept.Code)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	event, err := reg.GetEvent(scope.Code, concept.Code, "operation_failed")
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "operation_failed", event.Code)
	assert.Equal(t, "Operation Failed", event.Name)
	assert.Equal(t, "An operation on a task failed", event.Description)

}
