package user_handlers

import (
	"encoding/json"
	log "log/slog"
	"net/http"
	"strconv"
	"user-admin/internal/service"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserService *service.UserService
	Router *chi.Mux
}

func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		log.Error("Error getting users: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Error("Error encoding JSON: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
}

func (h *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.GetUserByID(int32(id))
	if err != nil {
		log.Error("Error retrieving user: ", err)
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    err = h.UserService.DeleteUser(int32(id))
    if err != nil {
        log.Error("Error deleting user: ", err)
        http.Error(w, "Error deleting user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User deleted successfully"))
}
