package server

import (
	"context"
	"slices"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// Authorize checks if the caller is an authorized Human (via Role) OR an authorized Machine (via Scope)
func Authorize(requiredRole, requiredScope string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			// 1. Extract claims from Kratos Context
			claimsInterface, ok := jwt.FromContext(ctx)
			if !ok {
				return nil, errors.Unauthorized("UNAUTHORIZED", "No token found in context")
			}

			claims, ok := claimsInterface.(jwtv5.MapClaims)
			if !ok {
				return nil, errors.Unauthorized("UNAUTHORIZED", "Invalid token claims format")
			}

			// 2. Edge Case: If no specific permissions are required, let them pass
			if requiredRole == "" && requiredScope == "" {
				return handler(ctx, req)
			}

			// 3. Check for Human Roles
			if requiredRole != "" {
				if role, roleOk := claims["role"].(string); roleOk && role == requiredRole {
					return handler(ctx, req) // Success: Human authorized
				}
			}

			// 4. Check for Machine Scopes
			if requiredScope != "" {
				// Format A: Space-separated string (Standard OAuth2)
				if scopeStr, scopeOk := claims["scope"].(string); scopeOk {
					if slices.Contains(strings.Split(scopeStr, " "), requiredScope) {
						return handler(ctx, req) // Success: Machine authorized
					}
				}
				// Format B: JSON Array of strings (Alternative)
				if scopeArr, scopeArrOk := claims["scope"].([]any); scopeArrOk {
					for _, s := range scopeArr {
						if strScope, isStr := s.(string); isStr && strScope == requiredScope {
							return handler(ctx, req) // Success: Machine authorized
						}
					}
				}
			}

			// 5. If we reach this point, they matched neither the role nor the scope
			return nil, errors.Forbidden("FORBIDDEN", "Insufficient permissions")
		}
	}
}

// NewWhiteListMatcher creates a matcher that skips middleware for specific routes
func NewWhiteListMatcher(whiteList []string) selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		for _, path := range whiteList {
			if operation == path || strings.HasPrefix(operation, path) {
				return false // skip middleware
			}
		}
		return true // execute middleware
	}
}
