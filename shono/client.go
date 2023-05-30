package shono

type Client interface {
	ScopeRepo
	ResourceRepo

	Close()
}
