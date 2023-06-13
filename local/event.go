package local

import "github.com/shono-io/shono/graph"

type eventRepo struct {
	events map[string]graph.Event
}

func (e *eventRepo) GetEventByReference(reference graph.EventReference) (*graph.Event, error) {
	res, fnd := e.events[reference.String()]
	if !fnd {
		return nil, nil
	}

	return &res, nil
}

func (e *eventRepo) GetEvent(scopeCode, conceptCode, code string) (*graph.Event, error) {
	return e.GetEventByReference(graph.EventReference{
		ScopeCode:   scopeCode,
		ConceptCode: conceptCode,
		Code:        code,
	})
}

func (e *eventRepo) ListEventsForConcept(scopeCode, conceptCode string) ([]graph.Event, error) {
	var res []graph.Event
	for _, v := range e.events {
		if v.ScopeCode == scopeCode && v.ConceptCode == conceptCode {
			res = append(res, v)
		}
	}
	return res, nil
}
