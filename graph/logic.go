package graph

type Expression string

type Logic interface {
	Kind() string
}
