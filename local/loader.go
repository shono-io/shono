package local

import (
	"github.com/shono-io/shono/backbone"
	"github.com/shono-io/shono/graph"
	"github.com/shono-io/shono/storage"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"os"
)

type spec struct {
	Backbone backbone.Config           `yaml:"backbone"`
	Storages map[string]storage.Config `yaml:"storages"`
	Scopes   map[string]scopeSpec      `yaml:"scopes"`
}

type scopeSpec struct {
	graph.Scope `yaml:",inline"`
	Concepts    map[string]conceptSpec `yaml:"concepts"`
}

type conceptSpec struct {
	graph.Concept `yaml:",inline"`
	Events        map[string]eventSpec `yaml:"events"`
}

type eventSpec struct {
	graph.Event `yaml:",inline"`
}

func LoadBytes(b []byte, opts ...Opt) (graph.Registry, error) {
	capitalize := cases.Title(language.English)

	b = []byte(os.ExpandEnv(string(b)))

	var inventory spec
	if err := yaml.Unmarshal(b, &inventory); err != nil {
		return nil, err
	}

	var o []Opt

	// -- backbone configuration
	bb, err := backbone.NewBackbone(inventory.Backbone)
	if err != nil {
		return nil, err
	}
	o = append(o, WithBackbone(bb))

	// -- scopes
	for scopeCode, scope := range inventory.Scopes {
		scope.Code = scopeCode
		if scope.Name == "" {
			scope.Name = capitalize.String(scopeCode)
		}

		if scope.Description == "" {
			scope.Description = scope.Name
		}

		o = append(o, WithScope(scope.Scope))

		for conceptCode, concept := range scope.Concepts {
			concept.ScopeCode = scopeCode
			concept.Code = conceptCode

			if concept.Name == "" {
				concept.Name = capitalize.String(conceptCode)
			}

			if concept.Description == "" {
				concept.Description = concept.Name
			}

			o = append(o, WithConcept(concept.Concept))

			for eventCode, event := range concept.Events {
				event.ScopeCode = scopeCode
				event.ConceptCode = conceptCode
				event.Code = eventCode

				if event.Name == "" {
					event.Name = capitalize.String(eventCode)
				}

				if event.Description == "" {
					event.Description = event.Name
				}
				o = append(o, WithEvent(event.Event))
			}
		}
	}

	// -- storages
	for k, s := range inventory.Storages {
		st, err := storage.NewStorage(k, s)
		if err != nil {
			return nil, err
		}

		o = append(o, WithStorage(st))
	}

	for _, opt := range opts {
		o = append(o, opt)
	}

	return NewRegistry(o...), nil
}

func Load(fs fs.FS, path string, opts ...Opt) (graph.Registry, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return LoadBytes(b, opts...)
}
