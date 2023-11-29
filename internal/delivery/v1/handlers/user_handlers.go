package user_handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"user-admin/internal/domain"
	"user-admin/internal/service"
	"user-admin/pkg/lib/utils"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserService *service.UserService
	Router *chi.Mux
}

func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page <= 0 {
        page = 1
    }

    pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
    if err != nil || pageSize <= 0 {
        pageSize = 12 // Default page size
    }

    users, err := h.UserService.GetAllUsers(page, pageSize)
    if err != nil {
        slog.Error("Error getting users: ", utils.Err(err))
        respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
        return
    }

    // Calculate previous and next pages
    previousPage := page - 1
    if previousPage < 1 {
        previousPage = 1
    }

    nextPage := page + 1

    response := struct {
        Users       *domain.UsersList `json:"users"`
        CurrentPage int               `json:"currentPage"`
        PrevPage    int               `json:"previousPage"`
        NextPage    int               `json:"nextPage"`
    }{
        Users:       users,
        CurrentPage: page,
        PrevPage:    previousPage,
        NextPage:    nextPage,
    }

    respondWithJSON(w, http.StatusOK, response)
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
		slog.Error("Error retrieving user: ", utils.Err(err))
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequest domain.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !utils.IsValidPhoneNumber(createUserRequest.PhoneNumber) {
		http.Error(w, "Invalid phone number format", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.CreateUser(&createUserRequest)
	if err != nil {
		slog.Error("Error creating user: ", utils.Err(err))
		http.Error(w, fmt.Sprintf("error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateUserRequest domain.UpdateUserRequest

	err = json.NewDecoder(r.Body).Decode(&updateUserRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateUserRequest.ID = int32(id)

	user, err := h.UserService.UpdateUser(&updateUserRequest)
	if err != nil {
		slog.Error("Error updating user: ", utils.Err(err))
		http.Error(w, fmt.Sprintf("error updating user: %v", err), http.StatusInternalServerError)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
        slog.Error("Error deleting user: ", utils.Err(err))
        http.Error(w, fmt.Sprintf("error deleting user: %s", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User deleted successfully"))
}

func (h *UserHandler) BlockUserHandler(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    err = h.UserService.BlockUser(int32(id))
    if err != nil {
        slog.Error("Error blocking user: ", utils.Err(err))
        http.Error(w, fmt.Sprintf("error blocking user: %s", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User blocked successfully"))
}

func (h *UserHandler) UnblockUserHandler(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    err = h.UserService.UnblockUser(int32(id))
    if err != nil {
        slog.Error("Error unblocking user by ID: ", utils.Err(err))
        http.Error(w, fmt.Sprintf("error unblocking user: %s", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User unblocked successfully"))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    w.WriteHeader(code)
    w.Write([]byte(message))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    if err := json.NewEncoder(w).Encode(payload); err != nil {
        slog.Error("Error encoding JSON: ", utils.Err(err))
        respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
    }
}