package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"user-admin/internal/config"
	"user-admin/pkg/lib/errors"
	"user-admin/pkg/lib/status"
	"user-admin/pkg/lib/utils"

	"log/slog"

	"github.com/dgrijalva/jwt-go"
)

// contextKey is a custom type for the context key used to store the JWT claims.
type contextKey string

const (
	// tokenKey is the context key for storing the JWT claims in the context.
	tokenKey contextKey = "token"
)

func AuthorizationMiddleware(cfg *config.Config, requiredRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractTokenFromHeader(r)
			if tokenString == "" {
				utils.RespondWithErrorJSON(w, status.Unauthorized, errors.AuthorizationTokenNotProvided)
				return
			}

			claims, err := validateToken(tokenString, cfg, isRefreshToken(r))
			if err != nil {
				utils.RespondWithErrorJSON(w, status.Unauthorized, fmt.Sprintf("Invalid authorization token: %v", err))
				return
			}

			adminRole, ok := claims["role"].(string)
			if !ok {
				utils.RespondWithErrorJSON(w, status.Unauthorized, errors.RoleNotFoundInTokenClaims)
				return
			}

			if !hasRequiredRole(adminRole, requiredRoles) {
				utils.RespondWithErrorJSON(w, status.Forbidden, errors.InsufficientPermission)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, tokenKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// validateToken validates the provided token string using the appropriate secret key based on token type.
func validateToken(tokenString string, cfg *config.Config, isRefreshToken bool) (jwt.MapClaims, error) {
	var secretKey string

	if isRefreshToken {
		secretKey = cfg.JWT.RefreshSecretKey
	} else {
		secretKey = cfg.JWT.AccessSecretKey
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
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

// isRefreshToken checks if the request is for a refresh token.
func isRefreshToken(r *http.Request) bool {
	return strings.Contains(r.URL.Path, "/refresh")
}

func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		slog.Error("Authorization header not found")
		return ""
	}

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
