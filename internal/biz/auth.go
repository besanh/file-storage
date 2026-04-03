package biz

import (
	"context"
	"crypto/rsa"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// PrivatePEM and PublicPEM are named types so wire can distinguish them.
type PrivatePEM []byte
type PublicPEM []byte

type AuthRepo interface {
	CheckPermission(ctx context.Context, subjectType, subjectID, relation, objectType, objectID string) (bool, error)
}

func NewPublicKey(pem PublicPEM) (*rsa.PublicKey, error) {
	return jwtv5.ParseRSAPublicKeyFromPEM(pem)
}

// GetUserID extracts the "sub" claim (User ID) from the context
func GetUserID(ctx context.Context) string {
	if claims, ok := jwt.FromContext(ctx); ok {
		if mapClaims, ok := claims.(jwtv5.MapClaims); ok {
			if sub, ok := mapClaims["sub"].(string); ok {
				return sub
			}
		}
	}
	return ""
}

// GetRole extracts the "role" claim from the context
func GetRole(ctx context.Context) string {
	if claims, ok := jwt.FromContext(ctx); ok {
		if mapClaims, ok := claims.(jwtv5.MapClaims); ok {
			if role, ok := mapClaims["role"].(string); ok {
				return role
			}
		}
	}
	return ""
}
