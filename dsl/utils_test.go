package dsl

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"strings"
	"testing"
)

func TestCleanMultilineStringWhitespace(t *testing.T) {
	input := strings.TrimSpace(`
reactor:
  code: on_sync_requested
  summary: Handle the SyncRequested event
  for:
    scope: govuk
    code: company
  when:
    scope: govuk
    concept: company
    code: sync_requested
  then:
    - label: "download_data_dump"
      branch:
        steps:
          - transform: |-
              root = if this.length() > 0 {
                "https://download.companieshouse.gov.uk/" + this.index(0)
              } else { 
                throw("url not found inside page")
              }
`)

	expected := strings.TrimSpace(`
root = if this.length() > 0 {
  "https://download.companieshouse.gov.uk/" + this.index(0)
} else { 
  throw("url not found inside page")
}`)

	var data map[string]any
	if err := yaml.Unmarshal([]byte(input), &data); err != nil {
		t.Fatalf("failed to unmarshal input: %v", err)
	}

	t.Log("data: ", data)

	assert.Equal(t, expected, data["reactor"].(map[string]any)["then"].([]any)[0].(map[string]any)["branch"].(map[string]any)["steps"].([]any)[0].(map[string]any)["transform"])
}
