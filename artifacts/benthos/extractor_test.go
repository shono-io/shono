package benthos

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/local"
	"github.com/shono-io/shono/system/file"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractorGenerator(t *testing.T) {
	shouldGenerateAValidExtractorArtifact(t)
}

func baseExtractor() *inventory.ExtractorBuilder {
	return inventory.NewExtractor("test", "my_extractor").
		InputEvent("test", "my_concept", "my_event").
		Output(file.NewOutput("out_file", file.WithOutputPath("/tmp/out.txt"))).
		Logic(inventory.NewLogic().Steps(dsl.Log("INFO").Message("my message")))
}

func shouldGenerateAValidExtractorArtifact(t *testing.T) {
	t.Run("should generate a valid extractor artifact", func(t *testing.T) {
		ext := baseExtractor().Build()

		inv, err := local.NewInventory().
			Extractor(ext).
			Event(inventory.NewEvent("test", "my_concept", "my_event").Build()).
			Build()
		if err != nil {
			t.Fatal(err)
		}

		art, err := NewExtractorGenerator().Generate("test", "my_extractor", inv, ext.Reference())
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, art)
		assert.Equal(t, "scopes/test/artifacts/my_extractor", art.Ref.String())
		assert.Equal(t, commons.ArtifactTypeExtractor, art.Type)
		assert.NotNil(t, art.Input)
		assert.NotNil(t, art.Output)
		assert.NotNil(t, art.DLQ)

		assert.Nil(t, art.Concept)

		// -- logic
		assert.NotNil(t, art.Logic)
		assert.NotEmpty(t, art.Logic.Steps)
		assert.NotEmpty(t, art.InputEvents)
		assert.Emptyf(t, art.OutputEvents, "extractors should not have output events")
	})
}
