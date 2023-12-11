package handlers

import (
	"encoding/json"
	"net/http"
	"user-admin/internal/domain"
	"user-admin/internal/service"

	"github.com/go-chi/chi/v5"
)

type AdminAuthHandler struct {
	AdminAuthService service.AdminAuthService
	Router *chi.Mux
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// SuccessResponse represents the structure of a success response.
type SuccessResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func (h *AdminAuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	token, err := h.AdminAuthService.LoginAdmin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		if err == domain.ErrAdminNotFound {
			// Handle case where admin/user was not found
			Error(w, http.StatusNotFound, "User not found")
			return
		}

		// Handle other errors
		Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	Success(w, http.StatusOK, LoginResponse{Token: token})
}

func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errResponse := ErrorResponse{Status: status, Message: message}
	json.NewEncoder(w).Encode(errResponse)
}

func Success(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	successResponse := SuccessResponse{Status: status, Data: data}
	json.NewEncoder(w).Encode(successResponse)
}