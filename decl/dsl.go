package decl

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/inventory"
)

type StepSpec struct {
	AddToStore      *AddToStoreStepSpec      `yaml:"addToStore,omitempty"`
	AsSuccessEvent  *AsSuccessEventStepSpec  `yaml:"asSuccessEvent,omitempty"`
	AsFailureEvent  *AsFailureEventStepSpec  `yaml:"asFailureEvent,omitempty"`
	Catch           *CatchStepSpec           `yaml:"catch,omitempty"`
	GetFromStore    *GetFromStoreStepSpec    `yaml:"getFromStore,omitempty"`
	ListFromStore   *ListFromStoreStepSpec   `yaml:"listFromStore,omitempty"`
	Log             *LogStepSpec             `yaml:"log,omitempty"`
	Raw             *RawStepSpec             `yaml:"raw,omitempty"`
	RemoveFromStore *RemoveFromStoreStepSpec `yaml:"removeFromStore,omitempty"`
	SetInStore      *SetInStoreStepSpec      `yaml:"setInStore,omitempty"`
	Switch          *SwitchStepSpec          `yaml:"switch,omitempty"`
}

func (ss *StepSpec) Children() []StepSpec {
	if ss.Catch != nil {
		return ss.Catch.Children()
	}
	if ss.Switch != nil {
		return ss.Switch.Children()
	}

	return nil
}

func (ss *StepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	if ss.AddToStore != nil {
		return ss.AddToStore.AsLogicStep(parent)
	}

	if ss.AsSuccessEvent != nil {
		return ss.AsSuccessEvent.AsLogicStep(parent)
	}

	if ss.AsFailureEvent != nil {
		return ss.AsFailureEvent.AsLogicStep(parent)
	}

	if ss.Catch != nil {
		return ss.Catch.AsLogicStep(parent)
	}

	if ss.GetFromStore != nil {
		return ss.GetFromStore.AsLogicStep(parent)
	}

	if ss.ListFromStore != nil {
		return ss.ListFromStore.AsLogicStep(parent)
	}

	if ss.Log != nil {
		return ss.Log.AsLogicStep(parent)
	}

	if ss.Raw != nil {
		return ss.Raw.AsLogicStep(parent)
	}

	if ss.RemoveFromStore != nil {
		return ss.RemoveFromStore.AsLogicStep(parent)
	}

	if ss.SetInStore != nil {
		return ss.SetInStore.AsLogicStep(parent)
	}

	if ss.Switch != nil {
		return ss.Switch.AsLogicStep(parent)
	}

	panic("incompatible logic step")
}

type StepSpecWithChildren interface {
	Children() []StepSpec
}

type AddToStoreStepSpec struct {
	Concept ConceptRef `yaml:"concept"`
	Key     string     `yaml:"key"`
}

func (a AddToStoreStepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	return dsl.AddToStore(a.Concept.Scope, a.Concept.Code, a.Key)
}

type AsSuccessEventStepSpec struct {
	EventCode string `yaml:"event"`
	Status    int    `yaml:"status"`
}

func (s AsSuccessEventStepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	return dsl.AsSuccessEvent(parent.Child("events", s.EventCode), s.Status, "this")
}

type AsFailureEventStepSpec struct {
	EventCode string `yaml:"event"`
	ErrorCode int    `yaml:"status"`
	Reason    string `yaml:"reason"`
}

func (s AsFailureEventStepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	return dsl.AsFailedEvent(parent.Child("events", s.EventCode), s.ErrorCode, s.Reason)
}

type CatchStepSpec []StepSpec

func (c CatchStepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	var steps []inventory.LogicStep
	for _, step := range c {
		steps = append(steps, step.AsLogicStep(parent))
	}

	return dsl.Catch(steps...)
}

func (c CatchStepSpec) Children() []StepSpec { return c }

type GetFromStoreStepSpec struct {
	Concept ConceptRef `yaml:"concept"`
	Key     string     `yaml:"key"`
}

func (g GetFromStoreStepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	return dsl.GetFromStore(g.Concept.Scope, g.Concept.Code, g.Key)
}

type ListFromStoreStepSpec struct {
	Concept ConceptRef `yaml:"concept"`
	Filter  []Filter   `yaml:"filters,omitempty"`
}

func (l ListFromStoreStepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	filters := map[string]string{}
	for _, filter := range l.Filter {
		filters[filter.Field] = filter.Value
	}
	return dsl.ListFromStore(l.Concept.Scope, l.Concept.Code, filters)
}

type LogStepSpec struct {
	Level   string `yaml:"level"`
	Message string `yaml:"message"`
}

func (l LogStepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	return dsl.Log(dsl.LogLevel(l.Level), l.Message)
}

type RawStepSpec map[string]any

func (r RawStepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	return dsl.Raw(r)
}

type RemoveFromStoreStepSpec struct {
	Concept ConceptRef `yaml:"concept"`
	Key     string     `yaml:"key"`
}

func (r RemoveFromStoreStepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	return dsl.RemoveFromStore(r.Concept.Scope, r.Concept.Code, r.Key)
}

type SetInStoreStepSpec struct {
	Concept ConceptRef `yaml:"concept"`
	Key     string     `yaml:"key"`
}

func (s SetInStoreStepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	return dsl.SetInStore(s.Concept.Scope, s.Concept.Code, s.Key)
}

type SwitchStepSpec []SwitchCase

func (s SwitchStepSpec) AsLogicStep(parent commons.Reference) inventory.LogicStep {
	cases := make([]dsl.ConditionalCase, 0, len(s))
	for _, c := range s {
		var steps []inventory.LogicStep
		for _, step := range c.Steps {
			steps = append(steps, step.AsLogicStep(parent))
		}

		cases = append(cases, dsl.ConditionalCase{
			Check: c.Condition,
			Steps: steps,
		})
	}
	return dsl.Switch(cases...)
}

func (s SwitchStepSpec) Children() []StepSpec {
	children := make([]StepSpec, 0, len(s))
	for _, c := range s {
		children = append(children, c.Steps...)
	}
	return children
}

type SwitchCase struct {
	Condition string     `yaml:"condition,omitempty"`
	Steps     []StepSpec `yaml:"steps"`
}

type ConceptRef struct {
	Scope string `yaml:"scope"`
	Code  string `yaml:"code"`
}

type EventRef struct {
	Scope   string `yaml:"scope"`
	Concept string `yaml:"concept"`
	Code    string `yaml:"code"`
}

type Filter struct {
	Field string `yaml:"field"`
	Value string `yaml:"value"`
}
