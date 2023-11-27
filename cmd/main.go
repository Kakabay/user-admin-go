package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"user-admin/internal/config"
	handlers "user-admin/internal/delivery/v1/handlers"
	repository "user-admin/internal/repository/postgres"
	"user-admin/internal/service"
	database "user-admin/pkg/database"
	utils "user-admin/pkg/lib/utils"
	"user-admin/pkg/logger"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.SetupLogger(cfg.Env)

	slog.Info("Starting the server...", slog.String("env", cfg.Env))
	slog.Debug("Debug messages are enabled") // If env is set to prod, debug messages are going to be disabled

	db, err := database.InitDB(cfg)
	if err != nil {
		slog.Error("Failed to init database:", utils.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	userRepository := repository.NewPostgresUserRepository(db.GetDB())
	userService := service.NewUserService(userRepository)
	userHandler := handlers.UserHandler{
		UserService: userService, 
		Router: chi.NewRouter(),
	}

	userHandler.Router.Get("/user", userHandler.GetAllUsersHandler)
	userHandler.Router.Get("/user/{id}", userHandler.GetUserByIDHandler)
	userHandler.Router.Post("/user", userHandler.CreateUserHandler)
	userHandler.Router.Post("/user/{id}", userHandler.UpdateUserHandler)
	userHandler.Router.Delete("/user/{id}", userHandler.DeleteUserHandler)
	userHandler.Router.Post("/user/{id}/block", userHandler.BlockUserHandler)
	userHandler.Router.Post("/user/{id}/unblock", userHandler.UnblockUserHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Info("Shutting down the server gracefully...")

		if err := db.Close(); err != nil {
			slog.Error("Error closing database:", utils.Err(err))
		}
		os.Exit(0)
	}()

	err = http.ListenAndServe(":8082", userHandler.Router)
	if err != nil {
		slog.Error("Server failed to start:", utils.Err(err))
	}
}
