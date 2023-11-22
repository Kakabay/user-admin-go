package user_handlers

import (
	"encoding/json"
	log "log/slog"
	"net/http"
	"user-admin/internal/service"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) NewRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger) // it is optional, but it let's you to log middleware itself
	router.Use(middleware.Recoverer)

	router.Get("/users", h.GetAllUsersHandler)

	return router
}

func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		log.Error("Error getting users:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Error("Error encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
}