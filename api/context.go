package api

import (
	"context"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// CustomClaims contains custom data we want from the token.
type ShonoClaims struct {
	Scope string `json:"scope"`
	OrgId string `json:"https://api.shono.io/org"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c ShonoClaims) Validate(ctx context.Context) error {
	return nil
}

func SubjectFromContext(ctx context.Context) *string {
	claims := Claims(ctx)
	if claims == nil {
		return nil
	}
	return &claims.RegisteredClaims.Subject
}

func Claims(ctx context.Context) *validator.ValidatedClaims {
	return ctx.Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
}

func CustomClaims(ctx context.Context) *ShonoClaims {
	claims := Claims(ctx)
	if claims == nil || claims.CustomClaims == nil {
		return nil
	}
	return claims.CustomClaims.(*ShonoClaims)
}

//
//func RoleFromContext(ctx context.Context) *auth.Role {
//	val := ctx.Value("role")
//	if val == nil {
//		return nil
//	}
//	return val.(*auth.Role)
//}

//func User(ctx *gin.Scope, g *graph.Graph) *graph.User {
//	claims := Claims(ctx)
//	if claims == nil {
//		return nil
//	}
//
//	user, err := g.User(ctx, claims.RegisteredClaims.Subject)
//	if err != nil {
//		logrus.Warn("failed to get user from context: ", err)
//		return nil
//	}
//
//	return user
//}
