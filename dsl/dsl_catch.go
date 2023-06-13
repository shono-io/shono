package dsl

import "github.com/shono-io/shono/graph"

func Catch(elements ...graph.Logic) CatchLogic {
	return CatchLogic{
		Logics: elements,
	}
}

type CatchLogic struct {
	Logics []graph.Logic
}

func (e CatchLogic) Kind() string {
	return "catch"
}
