package graph

func Catch(elements ...Logic) CatchLogic {
	return CatchLogic{
		Logics: elements,
	}
}

type CatchLogic struct {
	Logics []Logic
}

func (e CatchLogic) Kind() string {
	return "catch"
}
