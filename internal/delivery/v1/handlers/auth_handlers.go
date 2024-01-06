package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"user-admin/internal/domain"
	"user-admin/internal/service"
	"user-admin/pkg/lib/errors"
	"user-admin/pkg/lib/status"
	"user-admin/pkg/lib/utils"

	"github.com/go-chi/chi/v5"
)

type AdminAuthHandler struct {
	AdminAuthService service.AdminAuthService
	Router           *chi.Mux
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (h *AdminAuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		slog.Error("Error decoding login request:", utils.Err(err))
		Error(w, status.BadRequest, errors.InvalidRequestFormat)
		return
	}

	accessToken, refreshToken, err := h.AdminAuthService.LoginAdmin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		switch err {
		case domain.ErrAdminNotFound:
			Error(w, status.NotFound, errors.AdminNotFound)
		default:
			slog.Error("Error during login:", utils.Err(err))
			Error(w, status.Unauthorized, errors.InvalidCredentials)
		}
		return
	}

	loginResponse := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	respondJSON(w, status.OK, loginResponse)
}
func (h *AdminAuthHandler) RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := extractTokenFromHeader(r)
	if refreshToken == "" {
		Error(w, status.Unauthorized, errors.RefreshTokenNotProvided)
		return
	}

	newAccessToken, newRefreshToken, err := h.AdminAuthService.RefreshTokens(refreshToken)
	if err != nil {
		slog.Error("Error refreshing tokens:", utils.Err(err))
		Error(w, status.Unauthorized, errors.InvalidRefreshToken)
		return
	}

	respondJSON(w, status.OK, map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		slog.Error("Authorization header not found")
		return ""
	}

	return strings.TrimPrefix(bearerToken, "Bearer ")
}

func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errResponse := ErrorResponse{Status: status, Message: message}
	json.NewEncoder(w).Encode(errResponse)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
