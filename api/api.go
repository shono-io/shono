package api

type Opt func(*Api)

func WithRoute(route ...Route) Opt {
	return func(a *Api) {
		a.routes = append(a.routes, route...)
	}
}

func NewApi(opts ...Opt) *Api {
	a := &Api{
		routes: []Route{},
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

type Api struct {
	routes []Route
}
