package graph

type EventReference struct {
	ScopeCode   string `yaml:"scopeCode"`
	ConceptCode string `yaml:"conceptCode"`
	Code        string `yaml:"code"`
}

func (r EventReference) String() string {
	return r.ScopeCode + "__" + r.ConceptCode + "__" + r.Code
}

type Event struct {
	EventReference
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type EventRepo interface {
	GetEventByReference(reference EventReference) (*Event, error)
	GetEvent(scopeCode, conceptCode, code string) (*Event, error)
	ListEventsForConcept(scopeCode, conceptCode string) ([]Event, error)
}
