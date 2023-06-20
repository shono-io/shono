package local

import (
	"fmt"
	"github.com/hack-pad/hackpadfs"
	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/blobfs"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/hairyhenderson/go-fsimpl/httpfs"
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/commons"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"net/url"
	gos "os"
	"strings"
)

func LoadArtifact(uri string) (artifacts.Artifact, error) {
	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	mux.Add(blobfs.FS)
	mux.Add(gitfs.FS)
	mux.Add(httpfs.FS)

	fs, err := mux.Lookup(uri)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	p := strings.TrimSuffix(u.Path, ".yaml")
	ref, err := commons.ParseString(p)
	if err != nil {
		return nil, err
	}

	return (&ArtifactLoader{filesystem: fs}).LoadArtifact(ref)
}

func DumpArtifact(dir string, artifact artifacts.Artifact) error {
	return (&ArtifactDumper{}).StoreArtifact(artifact)
}

type ArtifactLoader struct {
	filesystem hackpadfs.FS
}

func (a *ArtifactLoader) LoadArtifact(ref commons.Reference) (artifacts.Artifact, error) {
	_, filename := referenceToFsName(ref)
	logrus.Tracef("loading artifact %q from %s", ref, filename)

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

	dir, filename := referenceToFsName(artifact.Reference())

	// -- make sure the directory exists
	if err := gos.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// -- create and write the file
	return gos.WriteFile(filename, b, 0644)
}

func referenceToFsName(ref commons.Reference) (dir string, filename string) {
	dir = fmt.Sprintf("%s/%s", ref.Parent().String(), ref.Kind())
	filename = fmt.Sprintf("%s.yaml", ref.Code())
	return dir, filename
}

func encodeArtifact(artifact artifacts.Artifact) ([]byte, error) {
	return yaml.Marshal(artifact)
}

func decodeArtifact(b []byte) (artifacts.Artifact, error) {
	var result artifacts.Artifact
	if err := yaml.Unmarshal(b, &result); err != nil {
		return nil, err
	}

	return result, nil
}
