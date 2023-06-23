package local

import (
	"fmt"
	"github.com/hack-pad/hackpadfs"
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/artifacts/benthos"
	"github.com/shono-io/shono/commons"
	"gopkg.in/yaml.v3"
	"io"
	gos "os"
)

func LoadArtifact(uri string) (artifacts.Artifact, error) {
	return (&ArtifactLoader{filesystem: gos.DirFS(".")}).LoadArtifact(uri)
}

func DumpArtifact(artifact artifacts.Artifact) error {
	return (&ArtifactDumper{}).StoreArtifact(artifact)
}

type ArtifactLoader struct {
	filesystem hackpadfs.FS
}

func (a *ArtifactLoader) LoadArtifact(filename string) (artifacts.Artifact, error) {
	f, err := a.filesystem.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// -- read the contents of the file
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return decodeArtifact(b)
}

type ArtifactDumper struct {
}

func (a *ArtifactDumper) StoreArtifact(artifact artifacts.Artifact) error {
	b, err := encodeArtifact(artifact)
	if err != nil {
		return err
	}

	filename := referenceToFsName(artifact.Reference())

	// -- create and write the file
	return gos.WriteFile(filename, b, 0644)
}

func referenceToFsName(ref commons.Reference) string {
	return fmt.Sprintf("%s.yaml", ref.Code())
}

func encodeArtifact(artifact artifacts.Artifact) ([]byte, error) {
	return yaml.Marshal(artifact)
}

func decodeArtifact(b []byte) (artifacts.Artifact, error) {
	var result benthos.Artifact
	if err := yaml.Unmarshal(b, &result); err != nil {
		panic(err)
	}

	return &result, nil
}