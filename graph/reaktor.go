package graph

import (
	"fmt"
)

func NewReaktorSpec(code string) *ReaktorSpec {
	return &ReaktorSpec{
		reaktor: &Reaktor{Code: code, Status: StatusExperimental},
	}
}

type ReaktorSpec struct {
	reaktor *Reaktor
}

func (r *ReaktorSpec) Summary(summary string) *ReaktorSpec {
	r.reaktor.Summary = summary
	return r
}

func (r *ReaktorSpec) Docs(docs string) *ReaktorSpec {
	r.reaktor.Docs = docs
	return r
}

func (r *ReaktorSpec) Stable() *ReaktorSpec {
	r.reaktor.Status = StatusStable
	return r
}

func (r *ReaktorSpec) Beta() *ReaktorSpec {
	r.reaktor.Status = StatusBeta
	return r
}

func (r *ReaktorSpec) Experimental() *ReaktorSpec {
	r.reaktor.Status = StatusExperimental
	return r
}

func (r *ReaktorSpec) Deprecated() *ReaktorSpec {
	r.reaktor.Status = StatusDeprecated
	return r
}

func (r *ReaktorSpec) On(scopeCode, conceptCode, eventCode string) *ReaktorSpec {

}

type Reaktor struct {
	Code    string `yaml:"code"`
	Status  Status `yaml:"status"`
	Summary string `yaml:"summary"`
	Docs    string `yaml:"docs"`

	Input   *EventReference `yaml:"input"`
	Logic   []Logic         `yaml:"logic"`
	Outputs []ReaktorOutput `yaml:"outputs"`
	Tests   []ReaktorTest   `yaml:"tests"`
}

type Trigger interface {
}

type ReaktorOutput struct {
	Event EventReference `yaml:"event"`
	Docs  string         `yaml:"docs"`
}

type ReaktorReference struct {
	ScopeCode   string `yaml:"scopeCode"`
	ConceptCode string `yaml:"conceptCode"`
	Code        string `yaml:"code"`
}

func (r ReaktorReference) String() string {
	return r.ScopeCode + "__" + r.ConceptCode + "__" + r.Code
}

type ReaktorBuilder struct {
	reaktor Reaktor
	outputs map[string]string
}

func (b *ReaktorBuilder) ExecuteFor(scopeCode, conceptCode string, logics ...Logic) *ReaktorBuilder {
	b.reaktor.ScopeCode = scopeCode
	b.reaktor.ConceptCode = conceptCode
	b.reaktor.Logic = logics
	return b
}

func (b *ReaktorBuilder) NamedAs(name string) *ReaktorBuilder {
	b.reaktor.Name = name
	return b
}

func (b *ReaktorBuilder) WithDocs(docs string) *ReaktorBuilder {
	b.reaktor.Docs = docs
	return b
}

func (b *ReaktorBuilder) Producing(eventCode, docs string) *ReaktorBuilder {
	b.outputs[eventCode] = docs

	return b
}

func (b *ReaktorBuilder) WithTest(test ReaktorTest) *ReaktorBuilder {
	b.reaktor.Tests = append(b.reaktor.Tests, test)
	return b
}

func (b *ReaktorBuilder) Build() (*Reaktor, error) {
	if b.reaktor.Input == nil {
		return nil, fmt.Errorf("no input event defined")
	}

	if b.reaktor.Logic == nil || len(b.reaktor.Logic) == 0 {
		return nil, fmt.Errorf("no logic defined")
	}

	if b.reaktor.ScopeCode == "" {
		return nil, fmt.Errorf("no scope code defined")
	}

	if b.reaktor.ConceptCode == "" {
		return nil, fmt.Errorf("no concept code defined")
	}

	if b.reaktor.Name == "" {
		b.reaktor.Name = fmt.Sprintf("On %s for concept %s in scope %s", b.reaktor.Input.Code, b.reaktor.ConceptCode, b.reaktor.ScopeCode)
	}

	for eventCode, docs := range b.outputs {
		b.reaktor.Outputs = append(b.reaktor.Outputs, ReaktorOutput{
			Event: EventReference{
				ScopeCode:   b.reaktor.ScopeCode,
				ConceptCode: b.reaktor.ConceptCode,
				Code:        eventCode,
			},
			Docs: docs,
		})
	}

	b.reaktor.Code = fmt.Sprintf("on(%s)", b.reaktor.Input.String())

	return &b.reaktor, nil
}

func InputEvent(scopeCode, conceptCode, eventCode string) *ReaktorBuilder {
	return &ReaktorBuilder{
		outputs: map[string]string{},
		reaktor: Reaktor{
			Input: &EventReference{
				ScopeCode:   scopeCode,
				ConceptCode: conceptCode,
				Code:        eventCode,
			},
		},
	}
}

type ReaktorRepo interface {
	GetReaktorByReference(reference ReaktorReference) (*Reaktor, error)
	GetReaktor(scopeCode, conceptCode, reaktorCode string) (*Reaktor, error)
	ListReaktorsForConcept(scopeCode, conceptCode string) ([]Reaktor, error)
}
