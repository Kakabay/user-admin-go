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

	authRouter := chi.NewRouter()
	mainRouter.Route("/auth", func(r chi.Router) {
		r.Mount("/", authRouter)
	})

	adminAuthRepository := repository.NewPostgresAdminAuthRepository(db.GetDB(), cfg.JWT)
	adminAuthService := service.NewAdminAuthService(adminAuthRepository)
	authHandler := handlers.AdminAuthHandler{
		AdminAuthService: *adminAuthService,
		Router:           authRouter,
	}

	authRouter.Post("/login", authHandler.LoginHandler)
	authRouter.Post("/refresh", authHandler.RefreshTokensHandler)

	userRouter := chi.NewRouter()
	userRouter.Use(middleware.AuthorizationMiddleware(cfg, []string{"admin", "super_admin"}))
	mainRouter.Route("/api/user", func(r chi.Router) {
		r.Mount("/", userRouter)
	})

	userRepository := repository.NewPostgresUserRepository(db.GetDB())
	userService := service.NewUserService(userRepository)
	userHandler := handlers.UserHandler{
		UserService: userService,
		Router:      userRouter,
	}

	userRouter.Get("/", userHandler.GetAllUsersHandler)
	userRouter.Get("/{id}", userHandler.GetUserByIDHandler)
	userRouter.Post("/", userHandler.CreateUserHandler)
	userRouter.Put("/{id}", userHandler.UpdateUserHandler)
	userRouter.Delete("/{id}", userHandler.DeleteUserHandler)
	userRouter.Post("/{id}/block", userHandler.BlockUserHandler)
	userRouter.Post("/{id}/unblock", userHandler.UnblockUserHandler)
	userRouter.Get("/search", userHandler.SearchUsersHandler)

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

	err = http.ListenAndServe(cfg.Address, mainRouter)
	if err != nil {
		slog.Error("Server failed to start:", utils.Err(err))
	}
}
