package runtime

import (
	"context"
	"fmt"
	_ "github.com/benthosdev/benthos/v4/public/components/all"
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/benthosdev/benthos/v4/public/service/servicetest"
	"github.com/shono-io/shono/benthos"
	"github.com/shono-io/shono/graph"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

func Test(env graph.Environment, logLevel string) (err error) {
	// -- generate the benthos configuration
	gen := benthos.NewGenerator()
	output, err := gen.Generate(context.Background(), env)
	if err != nil {
		return fmt.Errorf("failed to generate benthos configuration: %w", err)
	}

	if len(output.Streams) == 0 {
		logrus.Errorf("no streams to run")
		return nil
	}

	tmpDir := os.TempDir()
	if err := output.Dump(tmpDir); err != nil {
		return fmt.Errorf("failed to dump benthos configuration: %w", err)
	}

	// -- run the cli testing everything in the tmp dir
	servicetest.RunCLIWithArgs(context.Background(), "benthos", "test", "--log", logLevel, fmt.Sprintf("%s/*.yaml", tmpDir))

	return nil
}

func Run(env graph.Environment) (err error) {
	// -- generate the benthos configuration
	gen := benthos.NewGenerator()
	output, err := gen.Generate(context.Background(), env)
	if err != nil {
		return fmt.Errorf("failed to generate benthos configuration: %w", err)
	}

	if len(output.Streams) == 0 {
		logrus.Errorf("no streams to run")
		return nil
	}

	// -- create a waitgroup to wait for all the units to finish
	wg := sync.WaitGroup{}

	// -- each unit will become a benthos stream
	for _, stream := range output.Streams {
		go func(stream benthos.Stream) {
			wg.Add(1)
			logrus.Infof("starting stream %s", stream.Concept.Key())
			if err := runStreamLocally(stream); err != nil {
				logrus.Panicf("stream %s failed: %v", stream.Concept.Key(), err)
				wg.Done()
			}
		}(stream)
	}

	logrus.Infof("waiting for all streams to finish")

	wg.Wait()

	return nil
}

func runStreamLocally(stream benthos.Stream) error {
	// -- convert the unit to yaml
	b, err := yaml.Marshal(stream.Unit)
	if err != nil {
		return fmt.Errorf("failed to marshal unit to yaml: %w", err)
	}

	sb := service.NewStreamBuilder()
	if err := sb.SetYAML(string(b)); err != nil {
		return fmt.Errorf("failed to parse the configuration for the benthos stream: %w", err)
	}

	s, err := sb.Build()
	if err != nil {
		return fmt.Errorf("failed to build the benthos stream: %w", err)
	}

	return s.Run(context.Background())
}
