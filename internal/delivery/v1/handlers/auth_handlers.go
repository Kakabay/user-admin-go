package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"user-admin/internal/domain"
	"user-admin/internal/service"

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
		slog.Error("Error decoding login request:", err)
		Error(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	accessToken, refreshToken, err := h.AdminAuthService.LoginAdmin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		switch err {
		case domain.ErrAdminNotFound:
			Error(w, http.StatusNotFound, "User not found")
		default:
			slog.Error("Error during login:", err)
			Error(w, http.StatusUnauthorized, "Invalid credentials")
		}
		return
	}

	loginResponse := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	respondJSON(w, http.StatusOK, loginResponse)
}

func (h *AdminAuthHandler) RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Authorization")
	if refreshToken == "" {
		Error(w, http.StatusUnauthorized, "Refresh token not provided")
		return
	}

	newAccessToken, newRefreshToken, err := h.AdminAuthService.RefreshTokens(refreshToken)
	if err != nil {
		slog.Error("Error refreshing tokens:", err)
		Error(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	},
	)
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
