package shono

type Store interface {
	ConceptCode() string
	Entity

	AsBenthosComponent() (map[string]any, error)
}
