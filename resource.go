package go_shono

type Resource[T any] struct {
	Id            string
	ClientFactory func() T
}

func (r *Resource[T]) Client() T {
	return r.ClientFactory()
}
