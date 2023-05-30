package local

import (
	"github.com/shono-io/go-shono/shono"
)

func NewResourceRepo() shono.ResourceRepo {
	return &resourceRepo{
		resources: make(map[string]shono.Resource),
	}
}

type resourceRepo struct {
	resources map[string]shono.Resource
}

func (r *resourceRepo) GetResource(code string) (shono.Resource, bool, error) {
	res, fnd := r.resources[code]
	return res, fnd, nil
}

func (r *resourceRepo) AddResource(resource shono.Resource) error {
	r.resources[resource.GetCode()] = resource
	return nil
}

func (r *resourceRepo) RemoveResource(code string) error {
	delete(r.resources, code)
	return nil
}
