package graph

type Organization struct {
	Node
	Name string `json:"name"`
}

func (o *Organization) Domains() ([]Domain, error) {
	return nil, nil
}

func (o *Organization) Domain(key string) (*Domain, error) {
	return nil, nil
}
