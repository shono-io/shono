package go_shono

type RuntimeOpt func(*Runtime)

func WithResource(res ...Resource[any]) RuntimeOpt {
	return func(r *Runtime) {
		for _, resource := range res {
			r.resources[resource.Id] = resource
		}
	}
}

func NewRuntime(opt ...RuntimeOpt) *Runtime {
	r := &Runtime{
		resources: make(map[string]Resource[any]),
	}

	for _, o := range opt {
		o(r)
	}

	return r
}

type Runtime struct {
	resources map[string]Resource[any]
}

func (r *Runtime) Resource(id string) any {
	res, fnd := r.resources[id]
	if !fnd {
		return nil
	}

	return res.Client()
}
