package runtime

import (
	"context"
	"fmt"
	"github.com/benthosdev/benthos/v4/public/service/servicetest"
	"github.com/rs/xid"
	"github.com/shono-io/shono/artifacts"
	"github.com/shono-io/shono/storage"
	"github.com/sirupsen/logrus"
	"os"
)

type RunConfig struct {
	// the id of the application. This should not change between runs
	ApplicationId string `json:"application_id" yaml:"application_id"`

	StorageSystemId string `json:"storage" yaml:"storage"`
}

func configForArtifact(cfg RunConfig, systems SystemConfigs, artifact artifacts.Artifact, loglevel string) ([]byte, error) {
	if artifact == nil {
		return nil, fmt.Errorf("no artifact provided")
	}

	// -- configure the artifact input
	inp := artifact.Input()
	inpSystem, err := systems.Resolve(inp.Id)
	if err != nil {
		return nil, err
	}
	for k, v := range inpSystem.Config {
		inp.Config[k] = v
	}

	// -- configure the artifact output
	out := artifact.Output()
	outSystem, err := systems.Resolve(out.Id)
	if err != nil {
		return nil, err
	}
	for k, v := range outSystem.Config {
		out.Config[k] = v
	}

	// -- configure the artifact dlq
	dlq := artifact.Error()
	dlqSystem, err := systems.Resolve(dlq.Id)
	if err != nil {
		return nil, err
	}
	for k, v := range dlqSystem.Config {
		dlq.Config[k] = v
	}

	// -- generate the benthos configuration
	return GenerateBenthosConfig(artifact, loglevel)
}

func RunArtifact(cfg RunConfig, systems SystemConfigs, artifact artifacts.Artifact, loglevel string) error {
	ll, err := logrus.ParseLevel(loglevel)
	if err != nil {
		return fmt.Errorf("invalid log level %q: %w", loglevel, err)
	}
	logrus.SetLevel(ll)

	if artifact.Concept() != nil {
		// -- a concept might have a store associated with it
		if artifact.Concept().Stored {
			// -- find the storage system
			storageSystem, ok := systems[cfg.StorageSystemId]
			if !ok {
				return fmt.Errorf("storage system %q not found", cfg.StorageSystemId)
			}

			storage.Register(storageSystem.Kind, storageSystem.Config, false)
		}
	}

	// -- generate the benthos configuration
	benthosConfig, err := configForArtifact(cfg, systems, artifact, loglevel)
	if err != nil {
		return err
	}
	tmpFile := fmt.Sprintf("%s/%s.yaml", os.TempDir(), xid.New().String())
	if err := os.WriteFile(tmpFile, benthosConfig, 0644); err != nil {
		return fmt.Errorf("failed to write the artifact to temporary file %q: %w", tmpFile, err)
	}

	logrus.Infof("Running artifact %q from %q", artifact.Key(), tmpFile)
	servicetest.RunCLIWithArgs(context.Background(), "benthos", "-c", tmpFile, "--log.level", loglevel)

	//sb := service.NewStreamBuilder()
	//if err := sb.SetYAML(string(benthosConfig)); err != nil {
	//	return err
	//}
	//
	//sb.SetPrintLogger(logrus.StandardLogger())
	//
	//s, err := sb.Build()
	//if err != nil {
	//	return err
	//}
	//
	//return s.Run(context.Background())
	return nil
}

func TestArtifact(artifact artifacts.Artifact, loglevel string) error {
	// -- register the store
	storage.Register("memory", map[string]any{}, true)

	// -- generate the benthos configuration
	benthosConfig, err := GenerateBenthosConfig(artifact, loglevel)
	if err != nil {
		return err
	}

	tmpFile := fmt.Sprintf("%s/%s.yaml", os.TempDir(), xid.New().String())
	if err := os.WriteFile(tmpFile, benthosConfig, 0644); err != nil {
		return fmt.Errorf("failed to write the artifact to temporary file %q: %w", tmpFile, err)
	}

	servicetest.RunCLIWithArgs(context.Background(), "benthos", "test", "--log", loglevel, tmpFile)
	return nil
}
