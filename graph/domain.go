package graph

type Domain struct {
	Node
	Name string `json:"name"`
}

func (o *Organization) CreateDomain(name string) (*Domain, error) {
	return nil, nil
}
