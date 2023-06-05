package graph

type Generator interface {
	Generate(env Environment) error
}
