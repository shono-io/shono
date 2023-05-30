package shono

type Store interface {
	ScopeCode() string
	ConceptCode() string
	Entity

	AsBenthosComponent() (map[string]any, error)
}
