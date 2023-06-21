package runtime

import (
	"context"
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/storage"
)

type RunConfig struct {
	// the id of the application. This should not change between runs
	ApplicationId string `json:"application_id" yaml:"application_id"`

	// the configuration provided for the input, output and dlq
	Input  map[string]any `json:"input" yaml:"input"`
	Output map[string]any `json:"output" yaml:"output"`
	Dlq    map[string]any `json:"dlq" yaml:"dlq"`

	Storage StorageConfig `json:"storage" yaml:"storage"`
}

type StorageConfig struct {
	Name   string         `json:"name" yaml:"name"`
	Config map[string]any `json:"config" yaml:"config"`
}

func RunArtifact(cfg RunConfig, artifact artifacts.Artifact) error {
	// -- register the store
	if cfg.Storage.Name != "" {
		storage.Register(cfg.Storage.Name, cfg.Storage.Config, false)
	}

	// -- configure the artifact input
	inp := artifact.Input()
	if cfg.Input != nil {
		for k, v := range cfg.Input {
			inp.Config[k] = v
		}
	}

	// -- configure the artifact output
	out := artifact.Output()
	if cfg.Output != nil {
		for k, v := range cfg.Output {
			out.Config[k] = v
		}
	}

	// -- configure the artifact dlq
	dlq := artifact.Error()
	if cfg.Dlq != nil {
		for k, v := range cfg.Dlq {
			dlq.Config[k] = v
		}
	}

	// -- generate the benthos configuration
	benthosConfig, err := GenerateBenthosConfig(artifact)
	if err != nil {
		return err
	}

	sb := service.NewStreamBuilder()
	if err := sb.SetYAML(string(benthosConfig)); err != nil {
		return err
	}

	s, err := sb.Build()
	if err != nil {
		return err
	}

	return s.Run(context.Background())
}

//func Test(a *benthos.Artifact, loglevel string) error {
//	tmpFile := fmt.Sprintf("%s/%s.yaml", os.TempDir(), xid.New().String())
//	if err := os.WriteFile(tmpFile, a.Spec.Content, 0644); err != nil {
//		return fmt.Errorf("failed to write the artifact to temporary file %q: %w", tmpFile, err)
//	}
//
//	servicetest.RunCLIWithArgs(context.Background(), "benthos", "test", "--log", loglevel, tmpFile)
//	return nil
//}
