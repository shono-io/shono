package graph

func NewEventSpec(code string) *EventSpec {
	return &EventSpec{
		event: &Event{
			Code: code,
		},
	}
}

type EventSpec struct {
	event *Event
}

func (e *EventSpec) Summary(summary string) *EventSpec {
	e.event.Summary = summary
	return e
}

func (e *EventSpec) Docs(docs string) *EventSpec {
	e.event.Docs = docs
	return e
}

func (e *EventSpec) Stable() *EventSpec {
	e.event.Status = StatusStable
	return e
}

func (e *EventSpec) Beta() *EventSpec {
	e.event.Status = StatusBeta
	return e
}

func (e *EventSpec) Experimental() *EventSpec {
	e.event.Status = StatusExperimental
	return e
}

func (e *EventSpec) Deprecated() *EventSpec {
	e.event.Status = StatusDeprecated
	return e
}

type EventReference struct {
	ScopeCode   string `yaml:"scopeCode"`
	ConceptCode string `yaml:"conceptCode"`
	Code        string `yaml:"code"`
}

func (r EventReference) String() string {
	return r.ScopeCode + "__" + r.ConceptCode + "__" + r.Code
}

type Event struct {
	Code    string `yaml:"code"`
	Status  Status `yaml:"status"`
	Summary string `yaml:"summary,omitempty"`
	Docs    string `yaml:"docs,omitempty"`
}

type EventRepo interface {
	GetEventByReference(reference EventReference) (*Event, error)
	GetEvent(scopeCode, conceptCode, code string) (*Event, error)
	ListEventsForConcept(scopeCode, conceptCode string) ([]Event, error)
}
