package dsl

import (
	"fmt"
	"github.com/shono-io/shono/inventory"
)

type BranchBuilder struct {
	*BranchStep
}

func (bb *BranchBuilder) Label(label string) *BranchBuilder {
	bb.BranchStep.label = label

	return bb
}

func (bb *BranchBuilder) Pre(pre string) *BranchBuilder {
	bb.BranchStep.Pre = pre

	return bb
}

func (bb *BranchBuilder) Post(post string) *BranchBuilder {
	bb.BranchStep.Post = post

	return bb
}

func (bb *BranchBuilder) Steps(steps ...inventory.StepBuilder) *BranchBuilder {
	bb.BranchStep.Steps = inventory.BuildAllSteps(steps...)

	return bb
}

func (bb *BranchBuilder) Build() inventory.LogicStep {
	return bb.BranchStep
}

func Branch() *BranchBuilder {
	return &BranchBuilder{
		&BranchStep{},
	}
}

type BranchStep struct {
	label string
	Pre   string
	Steps []inventory.LogicStep
	Post  string
}

func (b *BranchStep) Label() string {
	return b.label
}

func (b *BranchStep) Kind() string {
	return "branch"
}

func (b *BranchStep) Validate() error {
	if b.Steps == nil || len(b.Steps) == 0 {
		return fmt.Errorf("at least one step is required for a branch")
	}

	return nil
}

func (b *BranchStep) MarshalBenthos(trace string) (map[string]any, error) {
	trace = fmt.Sprintf("%s/%s", trace, b.Kind())

	if err := b.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", trace, err)
	}

	// marshall the steps
	var procs []map[string]any
	for idx, step := range b.Steps {
		proc, err := step.MarshalBenthos(fmt.Sprintf("%s/steps[%d]", trace, idx))
		if err != nil {
			return nil, err
		}
		procs = append(procs, proc)
	}

	br := map[string]any{
		"processors": procs,
	}

	if b.Pre != "" {
		br["request_map"] = b.Pre
	}

	if b.Post != "" {
		br["result_map"] = b.Post
	}

	result := map[string]any{
		"branch": br,
	}

	if b.label != "" {
		result["label"] = b.label
	}

	return result, nil
}
