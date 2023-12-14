package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"user-admin/internal/config"
	handlers "user-admin/internal/delivery/v1/handlers"
	"user-admin/internal/delivery/v1/middleware"
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

	mainRouter := chi.NewRouter()

	adminAuthRepository := repository.NewPostgresAdminAuthRepository(db.GetDB(), cfg.JWT)
	adminAuthService := service.NewAdminAuthService(adminAuthRepository)
	authHandler := handlers.AdminAuthHandler{
		AdminAuthService: *adminAuthService,
		Router:           chi.NewRouter(),
	}

	authHandler.Router.Post("/login", authHandler.LoginHandler)

	var requiredRoles = []string{"admin", "super_admin"}

	mainRouter.Use(middleware.AuthorizationMiddleware(cfg, requiredRoles))

	// Group the "/api" routes
	mainRouter.Route("/api", func(r chi.Router) {
		userRepository := repository.NewPostgresUserRepository(db.GetDB())
		userService := service.NewUserService(userRepository)
		userHandler := handlers.UserHandler{
			UserService: userService,
			Router:      chi.NewRouter(),
		}

		userHandler.Router.Get("/", userHandler.GetAllUsersHandler)
		userHandler.Router.Get("/{id}", userHandler.GetUserByIDHandler)
		userHandler.Router.Post("/", userHandler.CreateUserHandler)
		userHandler.Router.Post("/{id}", userHandler.UpdateUserHandler)
		userHandler.Router.Delete("/{id}", userHandler.DeleteUserHandler)
		userHandler.Router.Post("/{id}/block", userHandler.BlockUserHandler)
		userHandler.Router.Post("/{id}/unblock", userHandler.UnblockUserHandler)
		userHandler.Router.Get("/search", userHandler.SearchUsersHandler)

		r.Mount("/user", userHandler.Router)
	})

	// Group the "/auth" routes
	mainRouter.Route("/auth", func(r chi.Router) {
		r.Mount("/", authHandler.Router)
	})

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

	err = http.ListenAndServe(":8082", mainRouter)
	if err != nil {
		slog.Error("Server failed to start:", utils.Err(err))
	}
}
