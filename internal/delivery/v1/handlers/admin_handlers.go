package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user-admin/internal/domain"
	"user-admin/internal/service"
	"user-admin/pkg/lib/status"
	"user-admin/pkg/lib/utils"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	AdminService *service.AdminService
	Router       chi.Router
}

func (h *AdminHandler) CreateAdminHandler(w http.ResponseWriter, r *http.Request) {
	var admin domain.CreateAdminRequest
	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		utils.RespondWithErrorJSON(w, status.BadRequest, "Invalid request body")
		return
	}

	if admin.Username == "" || admin.Password == "" || admin.Role == "" {
		utils.RespondWithErrorJSON(w, status.BadRequest, "Username, password, and role are required fields")
		return
	}

	createdAdmin, err := h.AdminService.CreateAdmin(&admin)
	if err != nil {
		switch err {
		case domain.ErrAdminAlreadyExists:
			utils.RespondWithErrorJSON(w, status.Conflict, "Admin with the same username already exists")
		default:
			utils.RespondWithErrorJSON(w, status.InternalServerError, fmt.Sprintf("Error creating admin: %v", err))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdAdmin)
}
