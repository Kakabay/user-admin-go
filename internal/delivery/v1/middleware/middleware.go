package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"user-admin/internal/config"
	utils "user-admin/pkg/lib/utils"

	"log/slog"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

const (
	tokenKey contextKey = "token"
)

func AuthorizationMiddleware(cfg *config.Config, requiredRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractTokenFromBody(r)
			if tokenString == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Authorization token not provided in the request body")
				return
			}

			claims, err := validateToken(tokenString, cfg)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid authorization token: %v", err))
				return
			}

			roles, ok := claims["roles"].([]interface{})
			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "Roles not found in token claims")
				return
			}

			if !hasRequiredRoles(roles, requiredRoles) {
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
func extractTokenFromBody(r *http.Request) string {
	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		slog.Error("Error decoding request body:", err)
		return ""
	}

	token, ok := requestBody[string(tokenKey)].(string)
	if !ok {
		slog.Error("Token not found in request body")
		return ""
	}

	return token
}

func hasRequiredRoles(userRoles []interface{}, requiredRoles []string) bool {
	for _, requiredRole := range requiredRoles {
		found := false
		for _, userRole := range userRoles {
			if userRole.(string) == requiredRole {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func validateToken(tokenString string, cfg *config.Config) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(cfg.JWT.SecretKey), nil
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
