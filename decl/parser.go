package decl

import "gopkg.in/yaml.v3"

func ParseBytes(b []byte) (*Spec, error) {
	return (&parser{}).Parse(b)
}

type parser struct {
}

func (p *parser) Parse(b []byte) (*Spec, error) {
	var res Spec
	if err := yaml.Unmarshal(b, &res); err != nil {
		return nil, err
	}

	if res.Scope != nil {
		applyScopeDefaults(res.Scope)
		if err := res.Scope.Validate(); err != nil {
			return nil, err
		}
	}

	if res.Concept != nil {
		applyConceptDefaults(res.Concept)
		if err := res.Concept.Validate(); err != nil {
			return nil, err
		}
	}

	if res.Event != nil {
		applyEventDefaults(res.Event)
		if err := res.Event.Validate(); err != nil {
			return nil, err
		}
	}

	return &res, nil
}
