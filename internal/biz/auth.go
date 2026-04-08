package biz

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// PrivatePEM and PublicPEM are named types so wire can distinguish them.
type PrivatePEM []byte
type PublicPEM []byte

type AuthRepo interface {
	CheckPermission(ctx context.Context, subjectType, subjectID, relation, objectType, objectID string) (bool, error)
	WriteRelationship(ctx context.Context, resourceType, resourceID, relation, subjectType, subjectID string) error
	DeleteRelationship(ctx context.Context, resourceType, resourceID, relation, subjectType, subjectID string) error
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

// GetActorInfo extracts the actor's type and ID from the context.
// It prioritizes the "SpiceDB" claim (used by M2M auth), and falls back to "sub" as a "user".
func GetActorInfo(ctx context.Context) (string, string, error) {
	if claims, ok := jwt.FromContext(ctx); ok {
		if mapClaims, ok := claims.(jwtv5.MapClaims); ok {
			// Check for SpiceDB claim first (handles different casing possibilities)
			var spiceDB map[string]any
			if claim, ok := mapClaims["SpiceDB"].(map[string]any); ok {
				spiceDB = claim
			} else if claim, ok := mapClaims["spiceDB"].(map[string]any); ok {
				spiceDB = claim
			} else if claim, ok := mapClaims["spicedb"].(map[string]any); ok {
				spiceDB = claim
			}

			if spiceDB != nil {
				var subType, subID string
				if t, ok := spiceDB["Type"].(string); ok {
					subType = t
				} else if t, ok := spiceDB["type"].(string); ok {
					subType = t
				}

				if id, ok := spiceDB["ID"].(string); ok {
					subID = id
				} else if id, ok := spiceDB["id"].(string); ok {
					subID = id
				}

				if subType != "" && subID != "" {
					// Map 'client' from token to 'service' identity in SpiceDB
					if subType == "client" {
						subType = "service"
					}
					return subType, subID, nil
				}
			}

			// Fallback to "sub" and default to "user" type
			if sub, ok := mapClaims["sub"].(string); ok && sub != "" {
				return "user", sub, nil
			}
		}
	}
	return "", "", fmt.Errorf("unauthorized: missing actor identity")
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
