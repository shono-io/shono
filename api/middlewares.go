package api

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/openapi"
	"net/http"
	"net/url"
	"time"
)

func auth0Middleware(oa *openapi.Collector, issuerURL *url.URL, audience string) func(http.Handler) http.Handler {
	oa.Reflector().SpecEns().ComponentsEns().SecuritySchemesEns().WithMapOfSecuritySchemeOrRefValuesItem("auth0", openapi3.SecuritySchemeOrRef{
		SecurityScheme: &openapi3.SecurityScheme{
			OAuth2SecurityScheme: &openapi3.OAuth2SecurityScheme{
				Flows: openapi3.OAuthFlows{
					ClientCredentials: &openapi3.ClientCredentialsFlow{
						TokenURL: issuerURL.String(),

						MapOfAnything: map[string]interface{}{
							"audience": audience,
						},
					},
				},
			},
		},
	})

	provider := jwks.NewCachingProvider(issuerURL, time.Duration(5*time.Minute))

	jwtValidator, _ := validator.New(provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &ShonoClaims{}
			},
		),
	)

	jwtMiddleware := jwtmiddleware.New(jwtValidator.ValidateToken)

	return jwtMiddleware.CheckJWT
}
