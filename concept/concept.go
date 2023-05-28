package concept

type ApiRoute struct {
	Path   string
	Method string
}

type Concept struct {
	Name      string
	apiRoutes []ApiRoute
	schema    string
}
