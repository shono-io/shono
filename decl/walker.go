package decl

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Handler interface {
	OnScope(scope *ScopeSpec) error
	OnConcept(concept *ConceptSpec) error
	OnEvent(event *EventSpec) error
	OnReactor(reactor *ReactorSpec) error
}

func Walk(path string, h Handler) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return walkDir(path, h)
	} else {
		return walkFile(path, h)
	}
}

func walkDir(path string, h Handler) error {
	items, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, item := range items {
		fn := fmt.Sprintf("%s/%s", path, item.Name())
		logrus.Debugf("processing %s", fn)

		if item.IsDir() {
			if err := walkDir(fn, h); err != nil {
				return err
			}
		} else {
			if err := walkFile(fn, h); err != nil {
				return err
			}
		}
	}

	return nil
}

func walkFile(path string, h Handler) error {
	if !strings.HasSuffix(path, ".yaml") {
		return nil
	}

	logrus.Debugf("processing %s", path)

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	spec, err := ParseBytes(b)
	if err != nil {
		return err
	}

	if spec.Scope != nil {
		if err := h.OnScope(spec.Scope); err != nil {
			return err
		}
	}

	if spec.Concept != nil {
		if err := h.OnConcept(spec.Concept); err != nil {
			return err
		}
	}

	if spec.Event != nil {
		if err := h.OnEvent(spec.Event); err != nil {
			return err
		}
	}

	if spec.Reactor != nil {
		if err := h.OnReactor(spec.Reactor); err != nil {
			return err
		}
	}

	return nil
}
