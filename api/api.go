package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	"github.com/swaggest/swgui/v4emb"
	"net/http"
	"net/url"
	"time"
)

type ApiOpt func(*web.Service, chi.Router)

func WithRoutes(routes func(r chi.Router)) ApiOpt {
	return func(s *web.Service, r chi.Router) {
		routes(r)
	}
}

func NewApi(issuerUrl string, audience string) (*API, error) {
	iu, err := url.Parse(issuerUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the issuer url: %v", err)
	}

	service := web.DefaultService()
	service.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logrus.StandardLogger()}))
	service.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	service.Wrap(
		middleware.NoCache,
		middleware.Timeout(30*time.Second),
	)

	service.OpenAPI.Info.Title = "Shono API"
	service.OpenAPI.Info.Version = "1.0.0"
	service.Docs("/docs", v4emb.New)

	auth0Scheme := openapi3.SecurityScheme{
		OAuth2SecurityScheme: &openapi3.OAuth2SecurityScheme{
			Flows: openapi3.OAuthFlows{
				ClientCredentials: &openapi3.ClientCredentialsFlow{
					TokenURL: iu.String() + "oauth/token",
					MapOfAnything: map[string]interface{}{
						"audience": audience,
					},
				},
			},
		},
	}

	secured := service.Group(func(r chi.Router) {
		r.Use(auth0Middleware(service.OpenAPICollector, iu, audience))
		r.Use(nethttp.SecurityMiddleware(service.OpenAPICollector, "auth0", auth0Scheme))
		//r.Use(ContextMiddleware(db))
	})

	service.Mount("/debug", middleware.Profiler())

	a := &API{
		service: service,
		secured: secured,
	}

	return a, nil
}

type API struct {
	service *web.Service
	secured chi.Router
}

func (a *API) SecuredRouter() chi.Router {
	return a.secured
}

func (a *API) Run() error {
	logrus.Info("Starting the api")
	return http.ListenAndServe("localhost:3002", a.service)
}
