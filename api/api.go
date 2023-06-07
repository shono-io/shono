package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/shono-io/shono/graph"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func NewApi(env graph.Environment) (*API, error) {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logrus.StandardLogger()}),
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
		middleware.NoCache,
		middleware.Timeout(30*time.Second),
	)

	r.Mount("/debug", middleware.Profiler())

	if err := registerRoutesForEnv(r, env); err != nil {
		return nil, err
	}

	return &API{r}, nil
}

type API struct {
	r chi.Router
}

func (a *API) Run() error {
	logrus.Info("Starting the api")
	return http.ListenAndServe("localhost:3002", a.r)
}

func registerRoutesForEnv(r chi.Router, env graph.Environment) error {
	scopes, err := env.ListScopes()
	if err != nil {
		return err
	}

	for _, scope := range scopes {
		if err := registerRoutesForScope(r, env, scope); err != nil {
			return err
		}
	}

	return nil
}

func registerRoutesForScope(r chi.Router, env graph.Environment, scope graph.Scope) error {
	logrus.Infof("Registering routes for scope %s", scope.Name())

	bb, err := env.GetBackbone()
	if err != nil {
		return err
	}

	bbc, err := bb.GetClient()
	if err != nil {
		return err
	}

	concepts, err := env.ListConceptsForScope(scope.Key())
	if err != nil {
		return err
	}

	for _, concept := range concepts {
		if concept.Requests() == nil {
			continue
		}

		for _, request := range concept.Requests() {
			h, err := NewRequestHandler(env, bbc, request)
			if err != nil {
				return err
			}

			var path string
			var method string
			switch request.Kind {
			case graph.ListOperationType:
				method = http.MethodGet
				path = fmt.Sprintf("/%s/%s", scope.Key().Code(), concept.Plural())
			case graph.GetOperationType:
				method = http.MethodGet
				path = fmt.Sprintf("/%s/%s/{id}", scope.Key().Code(), concept.Plural())
			case graph.CreateOperationType:
				method = http.MethodPost
				path = fmt.Sprintf("/%s/%s", scope.Key().Code(), concept.Plural())
			case graph.UpdateOperationType:
				method = http.MethodPut
				path = fmt.Sprintf("/%s/%s/{id}", scope.Key().Code(), concept.Plural())
			case graph.DeleteOperationType:
				method = http.MethodDelete
				path = fmt.Sprintf("/%s/%s/{id}", scope.Key().Code(), concept.Plural())
			}

			if method != "" {
				r.Method(method, path, h)
			}
		}
	}

	return nil
}
