package dsl

import "github.com/shono-io/shono/graph"

func If(check graph.Expression, block ...graph.Logic) ConditionalLogic {
	return ConditionalLogic{
		Cases: []CaseLogic{
			{
				Check:  check,
				Logics: block,
			},
		},
	}
}

func Switch(cases ...CaseLogic) ConditionalLogic {
	return ConditionalLogic{
		Cases: cases,
	}
}

type ConditionalLogic struct {
	Cases []CaseLogic `yaml:"cases"`
}

func (e ConditionalLogic) Kind() string {
	return "conditional"
}

func SwitchCase(check graph.Expression, block ...graph.Logic) CaseLogic {
	return CaseLogic{
		Check:  check,
		Logics: block,
	}
}

func SwitchDefault(block ...graph.Logic) CaseLogic {
	return CaseLogic{
		Logics: block,
	}
}

type CaseLogic struct {
	Check  graph.Expression
	Logics []graph.Logic
}

func (e CaseLogic) Kind() string {
	return "case"
}
