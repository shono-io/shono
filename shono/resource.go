package shono

type Resource interface {
	Entity
}

type ResourceRepo interface {
	GetResource(code string) (Resource, bool, error)
	AddResource(resource Resource) error
	RemoveResource(code string) error
}
