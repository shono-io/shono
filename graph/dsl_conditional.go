package graph

func If(check Expression, block ...Logic) ConditionalLogic {
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
	Cases []CaseLogic
}

func (e ConditionalLogic) Kind() string {
	return "conditional"
}

func SwitchCase(check Expression, block ...Logic) CaseLogic {
	return CaseLogic{
		Check:  check,
		Logics: block,
	}
}

func SwitchDefault(block ...Logic) CaseLogic {
	return CaseLogic{
		Logics: block,
	}
}

type CaseLogic struct {
	Check  Expression
	Logics []Logic
}

func (e CaseLogic) Kind() string {
	return "case"
}
