package api

import "net/http"

type Route interface {
	Path() string
	Name() string
	Title() string
	Description() string
	Tags() []string
	AuthenticationRequired() bool
	Handler() http.Handler
}
