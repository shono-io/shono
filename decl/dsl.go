package decl

import (
	"fmt"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/inventory"
	"github.com/sirupsen/logrus"
)

type Steps []StepSpec

type StepSpec struct {
	AddToStore      *AddToStoreStepSpec      `yaml:"addToStore,omitempty"`
	AsSuccessEvent  *AsSuccessEventStepSpec  `yaml:"asSuccessEvent,omitempty"`
	AsFailureEvent  *AsFailureEventStepSpec  `yaml:"asFailureEvent,omitempty"`
	Branch          *BranchStepSpec          `yaml:"branch,omitempty"`
	Catch           *CatchStepSpec           `yaml:"catch,omitempty"`
	GetFromStore    *GetFromStoreStepSpec    `yaml:"getFromStore,omitempty"`
	ListFromStore   *ListFromStoreStepSpec   `yaml:"listFromStore,omitempty"`
	Log             *LogStepSpec             `yaml:"log,omitempty"`
	Raw             *RawStepSpec             `yaml:"raw,omitempty"`
	RemoveFromStore *RemoveFromStoreStepSpec `yaml:"removeFromStore,omitempty"`
	SetInStore      *SetInStoreStepSpec      `yaml:"setInStore,omitempty"`
	Switch          *SwitchStepSpec          `yaml:"switch,omitempty"`
	Transform       *TransformStepSpec       `yaml:"transform,omitempty"`
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

func (ss *StepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	if ss.AddToStore != nil {
		return ss.AddToStore.AsLogicStep(fmt.Sprintf("%s/add_to_store", path), parent)
	}

	if ss.AsSuccessEvent != nil {
		return ss.AsSuccessEvent.AsLogicStep(fmt.Sprintf("%s/as_success_event", path), parent)
	}

	if ss.AsFailureEvent != nil {
		return ss.AsFailureEvent.AsLogicStep(fmt.Sprintf("%s/as_failure_event", path), parent)
	}

	if ss.Branch != nil {
		return ss.Branch.AsLogicStep(fmt.Sprintf("%s/branch", path), parent)
	}

	if ss.Catch != nil {
		return ss.Catch.AsLogicStep(fmt.Sprintf("%s/catch", path), parent)
	}

	if ss.GetFromStore != nil {
		return ss.GetFromStore.AsLogicStep(fmt.Sprintf("%s/get_from_store", path), parent)
	}

	if ss.ListFromStore != nil {
		return ss.ListFromStore.AsLogicStep(fmt.Sprintf("%s/list_from_store", path), parent)
	}

	if ss.Log != nil {
		return ss.Log.AsLogicStep(fmt.Sprintf("%s/log", path), parent)
	}

	if ss.Raw != nil {
		return ss.Raw.AsLogicStep(fmt.Sprintf("%s/raw", path), parent)
	}

	if ss.RemoveFromStore != nil {
		return ss.RemoveFromStore.AsLogicStep(fmt.Sprintf("%s/remove_from_store", path), parent)
	}

	if ss.SetInStore != nil {
		return ss.SetInStore.AsLogicStep(fmt.Sprintf("%s/set_in_store", path), parent)
	}

	if ss.Switch != nil {
		return ss.Switch.AsLogicStep(fmt.Sprintf("%s/switch", path), parent)
	}

	if ss.Transform != nil {
		return ss.Transform.AsLogicStep(fmt.Sprintf("%s/transform", path), parent)
	}

	logrus.Panicf("unknown step type at %s", path)
	return nil
}

type StepSpecWithChildren interface {
	Children() []StepSpec
}

type BranchStepSpec struct {
	Pre   string     `yaml:"pre"`
	Post  string     `yaml:"post"`
	Steps []StepSpec `yaml:"steps"`
}

func (s BranchStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	var st []inventory.StepBuilder
	for idx, step := range s.Steps {
		st = append(st, step.AsLogicStep(fmt.Sprintf("%s/steps[%d]", path, idx), parent))
	}
	return dsl.Branch().Pre(s.Pre).Post(s.Post).Steps(st...)
}

type AddToStoreStepSpec struct {
	Concept ConceptRef `yaml:"concept"`
	Key     string     `yaml:"key"`
}

func (a AddToStoreStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	return dsl.AddToStore(a.Concept.Scope, a.Concept.Code, a.Key)
}

type AsSuccessEventStepSpec struct {
	EventCode string `yaml:"event"`
	Status    int    `yaml:"status"`
}

func (s AsSuccessEventStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	return dsl.AsSuccessEvent(parent.Child("events", s.EventCode), s.Status, "this")
}

type AsFailureEventStepSpec struct {
	EventCode string `yaml:"event"`
	ErrorCode int    `yaml:"status"`
	Reason    string `yaml:"reason"`
}

func (s AsFailureEventStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	return dsl.AsFailedEvent(parent.Child("events", s.EventCode), s.ErrorCode, s.Reason)
}

type CatchStepSpec []StepSpec

func (c CatchStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	var steps []inventory.StepBuilder
	for idx, step := range c {
		st := step.AsLogicStep(fmt.Sprintf("%s/clause[%d]", path, idx), parent)

		steps = append(steps, st)
	}

	return dsl.Catch().Steps(steps...)
}

func (c CatchStepSpec) Children() []StepSpec { return c }

type GetFromStoreStepSpec struct {
	Concept ConceptRef `yaml:"concept"`
	Key     string     `yaml:"key"`
}

func (g GetFromStoreStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	return dsl.GetFromStore(g.Concept.Scope, g.Concept.Code, g.Key)
}

type ListFromStoreStepSpec struct {
	Concept ConceptRef `yaml:"concept"`
	Filter  []Filter   `yaml:"filters,omitempty"`
}

func (l ListFromStoreStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
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

func (l LogStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	return dsl.Log(dsl.LogLevel(l.Level)).Message(l.Message)
}

type RawStepSpec map[string]any

func (r RawStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	return dsl.Raw().Content(r)
}

type RemoveFromStoreStepSpec struct {
	Concept ConceptRef `yaml:"concept"`
	Key     string     `yaml:"key"`
}

func (r RemoveFromStoreStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	return dsl.RemoveFromStore(r.Concept.Scope, r.Concept.Code, r.Key)
}

type SetInStoreStepSpec struct {
	Concept ConceptRef `yaml:"concept"`
	Key     string     `yaml:"key"`
}

func (s SetInStoreStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	return dsl.SetInStore(s.Concept.Scope, s.Concept.Code, s.Key)
}

type SwitchStepSpec []SwitchCase

func (s SwitchStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	sw := dsl.Switch()

	for _, c := range s {
		var steps []inventory.StepBuilder
		for idx, step := range c.Steps {
			steps = append(steps, step.AsLogicStep(fmt.Sprintf("%s/case[%d]", path, idx), parent))
		}

		if c.Condition == "" {
			sw.Default(steps...)
		} else {
			sw.Case(c.Condition, steps...)
		}
	}

	return sw
}

func (s SwitchStepSpec) Children() []StepSpec {
	children := make([]StepSpec, 0, len(s))
	for _, c := range s {
		children = append(children, c.Steps...)
	}
	return children
}

type TransformStepSpec string

func (s TransformStepSpec) AsLogicStep(path string, parent commons.Reference) inventory.StepBuilder {
	return dsl.Transform().Mapping(string(s))
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
