package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"user-admin/internal/config"
	utils "user-admin/pkg/lib/utils"

	"log/slog"

	"github.com/dgrijalva/jwt-go"
)

// contextKey is a custom type for the context key used to store the JWT claims.
type contextKey string

const (
	// tokenKey is the context key for storing the JWT claims in the context.
	tokenKey contextKey = "token"
)

// AuthorizationMiddleware is a middleware function that performs authorization checks based on JWT tokens.
func AuthorizationMiddleware(cfg *config.Config, requiredRoles []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tokenString := extractTokenFromHeader(r)
            if tokenString == "" {
                utils.RespondWithError(w, http.StatusUnauthorized, "Authorization token not provided")
                return
            }

            claims, err := validateToken(tokenString, cfg)
            if err != nil {
                utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid authorization token: %v", err))
                return
            }

			adminRole, ok := claims["role"].(string)
			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "Role not found in token claims")
				return
			}
			
			if !hasRequiredRole(adminRole, requiredRoles) {
				utils.RespondWithError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

            // Pass the claims to the next handler
            ctx := r.Context()
            ctx = context.WithValue(ctx, tokenKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func extractTokenFromHeader(r *http.Request) string {
    bearerToken := r.Header.Get("Authorization")
    if bearerToken == "" {
        slog.Error("Authorization header not found")
        return ""
    }

    // Check if the Authorization header has the expected "Bearer " prefix
    if !strings.HasPrefix(bearerToken, "Bearer ") {
        slog.Error("Invalid Authorization header format")
        return ""
    }

    return strings.TrimPrefix(bearerToken, "Bearer ")
}


func hasRequiredRole(adminRole string, requiredRoles []string) bool {
	for _, requiredRole := range requiredRoles {
		if adminRole == requiredRole {
			return true
		}
	}
	return false
}


// validateToken validates and parses the JWT token.
func validateToken(tokenString string, cfg *config.Config) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(cfg.JWT.AccessSecretKey), nil
    })

    if err != nil || !token.Valid {
        slog.Error("Token validation error: %v", err)
        return nil, fmt.Errorf("token validation error: %v", err)
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || claims == nil {
        slog.Error("Invalid token claims")
        return nil, fmt.Errorf("invalid token claims")
    }

    return claims, nil
}
