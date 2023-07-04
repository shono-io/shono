package decl

import "fmt"

type Spec struct {
	Scope   *ScopeSpec   `yaml:"scope,omitempty"`
	Concept *ConceptSpec `yaml:"concept,omitempty"`
	Event   *EventSpec   `yaml:"event,omitempty"`
}

type ScopeSpec struct {
	Code    string `yaml:"code"`
	Summary string `yaml:"summary"`
	Status  string `yaml:"status,omitempty"`
	Docs    string `yaml:"docs,omitempty"`
}

func (s *ScopeSpec) Validate() error {
	if s.Code == "" {
		return fmt.Errorf("scope code is required")
	}

	if s.Summary == "" {
		return fmt.Errorf("scope summary is required")
	}

	if s.Status == "" {
		return fmt.Errorf("scope status is required")
	}

	return nil
}

func applyScopeDefaults(scope *ScopeSpec) {
	if scope.Status == "" {
		scope.Status = "experimental"
	}
}

type ConceptSpec struct {
	Scope   string `yaml:"scope"`
	Code    string `yaml:"code"`
	Summary string `yaml:"summary"`
	Status  string `yaml:"status,omitempty"`
	Docs    string `yaml:"docs,omitempty"`
	Stored  bool   `yaml:"stored,omitempty"`
}

func (c *ConceptSpec) Validate() error {
	if c.Scope == "" {
		return fmt.Errorf("concept scope is required")
	}

	if c.Code == "" {
		return fmt.Errorf("concept code is required")
	}

	if c.Summary == "" {
		return fmt.Errorf("concept summary is required")
	}

	if c.Status == "" {
		return fmt.Errorf("concept status is required")
	}

	return nil
}

func applyConceptDefaults(concept *ConceptSpec) {
	if concept.Status == "" {
		concept.Status = "experimental"
	}
}

type EventSpec struct {
	Scope   string `yaml:"scope"`
	Concept string `yaml:"concept"`
	Code    string `yaml:"code"`
	Summary string `yaml:"summary"`
	Status  string `yaml:"status,omitempty"`
	Docs    string `yaml:"docs,omitempty"`
}

func (e *EventSpec) Validate() error {
	if e.Scope == "" {
		return fmt.Errorf("event scope is required")
	}

	if e.Concept == "" {
		return fmt.Errorf("event concept is required")
	}

	if e.Code == "" {
		return fmt.Errorf("event code is required")
	}

	if e.Summary == "" {
		return fmt.Errorf("event summary is required")
	}

	if e.Status == "" {
		return fmt.Errorf("event status is required")
	}

	return nil
}

func applyEventDefaults(event *EventSpec) {
	if event.Status == "" {
		event.Status = "experimental"
	}
}
