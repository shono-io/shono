package runtime

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type LoadOpt func([]string)

func WithPath(path string) LoadOpt {
	return func(args []string) {
		args = append(args, path)
	}
}

type SystemConfigs map[string]SystemConfig

func (s SystemConfigs) Resolve(id string) (*SystemConfig, error) {
	if id == "" {
		panic("id cannot be empty")
	}

	cfg, ok := s[id]
	if !ok {
		return nil, fmt.Errorf("system %q not found", id)
	}

	return &cfg, nil
}

type SystemConfig struct {
	Kind   string         `yaml:"kind"`
	Config map[string]any `yaml:"config"`
}

func LoadSystems(opt ...LoadOpt) (SystemConfigs, error) {
	paths := []string{
		"./systems.yaml",
		"~/.shono/systems.yaml",
	}

	for _, o := range opt {
		o(paths)
	}

	for _, p := range paths {
		if strings.HasPrefix(p, "~/") {
			p = strings.Replace(p, "~/", os.Getenv("HOME")+"/", 1)
		}

		if _, err := os.Stat(p); err == nil {
			return load(p)
		}
	}

	return nil, nil
}

func load(path string) (SystemConfigs, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg SystemConfigs
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	for k, _ := range cfg {
		logrus.Infof("loaded system %q", k)
	}

	return cfg, nil
}
