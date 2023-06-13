package runtime

import (
	"context"
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/local"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCatchLogicToProcessor(t *testing.T) {
	registry := local.NewRegistry()

	l := dsl.Catch(dsl.Log("ERROR", "my error message"))
	m, err := toProcessor(context.Background(), registry, l)
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"catch": []map[string]any{
			{
				"log": map[string]any{
					"level":   "ERROR",
					"message": "my error message",
				},
			},
		},
	}, m)
}

func TestSwitchLogicToProcessor(t *testing.T) {
	registry := local.NewRegistry()
	l := dsl.Switch(
		dsl.SwitchCase("this.status == 'NEW'", dsl.Log("INFO", "New task created")),
		dsl.SwitchDefault(dsl.Log("ERROR", "my error message")))
	m, err := toProcessor(context.Background(), registry, l)
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"switch": []map[string]any{
			{
				"check": "this.status == 'NEW'",
				"processors": []map[string]any{
					{
						"log": map[string]any{
							"level":   "INFO",
							"message": "New task created",
						},
					},
				},
			},
			{
				"processors": []map[string]any{
					{
						"log": map[string]any{
							"level":   "ERROR",
							"message": "my error message",
						},
					},
				},
			},
		},
	}, m)
}

func TestCaseLogicToProcessor(t *testing.T) {
	registry := local.NewRegistry()

	l := dsl.SwitchCase("this.status == 'NEW'", dsl.Log("ERROR", "my error message"))
	m, err := toProcessor(context.Background(), registry, l)
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"check": "this.status == 'NEW'",
		"processors": []map[string]any{
			{
				"log": map[string]any{
					"level":   "ERROR",
					"message": "my error message",
				},
			},
		},
	}, m)
}

func TestDefaultCaseLogicToProcessor(t *testing.T) {
	registry := local.NewRegistry()

	l := dsl.SwitchDefault(dsl.Log("ERROR", "my error message"))
	m, err := toProcessor(context.Background(), registry, l)
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"processors": []map[string]any{
			{
				"log": map[string]any{
					"level":   "ERROR",
					"message": "my error message",
				},
			},
		},
	}, m)
}

func TestLogLogicToProcessor(t *testing.T) {
	registry := local.NewRegistry()

	l := dsl.Log("ERROR", "my error message")
	m, err := toProcessor(context.Background(), registry, l)
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"log": map[string]any{
			"level":   "ERROR",
			"message": "my error message",
		},
	}, m)
}

func TestTransformLogicToProcessor(t *testing.T) {
	registry := local.NewRegistry()

	l := dsl.Transform(
		dsl.MapMeta("io_shono_kind", dsl.AsEventReference("todo", "tasks", "created")),
		dsl.Map("status", dsl.AsConstant(201)),
		dsl.Map("task", dsl.ToBloblang("this")),
	)
	m, err := toProcessor(context.Background(), registry, l)
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"mapping": strings.TrimSpace(`
meta io_shono_kind = "todo__tasks__created"
root.status = 201
root.task = this
`)}, m)
}

func TestTransformRootLogicToProcessor(t *testing.T) {
	registry := local.NewRegistry()

	l := dsl.Transform(
		dsl.MapMeta("io_shono_kind", dsl.AsEventReference("todo", "tasks", "created")),
		dsl.MapRoot(),
	)
	m, err := toProcessor(context.Background(), registry, l)
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"mapping": strings.TrimSpace(`
meta io_shono_kind = "todo__tasks__created"
root = this
`)}, m)
}
