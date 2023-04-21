package states

func NewCatalog() (*Catalog, error) {
	result := &Catalog{
		states: make(map[Kind]*State),
	}

	return result, nil
}

type Catalog struct {
	states map[Kind]*State
}

func (e *Catalog) RegisterState(kind Kind, t any) error {
	e.states[kind] = &State{
		kind:      kind,
		valueType: t,
	}

	return nil
}

func (e *Catalog) State(stateKind Kind) (*State, bool) {
	state, ok := e.states[stateKind]
	return state, ok
}
