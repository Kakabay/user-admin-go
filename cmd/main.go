// main.go

package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"user-admin/internal/config"
	"user-admin/internal/delivery/v1/middleware"
	"user-admin/internal/delivery/v1/routers"
	repository "user-admin/internal/repository/postgres"
	"user-admin/internal/service"
	"user-admin/pkg/database"
	utils "user-admin/pkg/lib/utils"
	"user-admin/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// @title Admin Panel
// version 1.0
// description API Server for Admin Panel

// @host localhost:8081
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

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

	authMiddlewareForAdmin := middleware.AuthMiddleware(cfg, []string{"admin"})
	authMiddlewareForSuperAdmin := middleware.AuthMiddleware(cfg, []string{"super_admin"})

	// Admin routes
	adminRouter := chi.NewRouter()
	adminRouter.Use(authMiddlewareForSuperAdmin) // Apply auth middleware to admin routes
	mainRouter.Route("/api/admin", func(r chi.Router) {
		r.Mount("/", adminRouter)
	})

	adminRepository := repository.NewPostgresAdminRepository(db.GetDB())
	adminService := service.NewAdminService(adminRepository)
	routers.SetupAdminRoutes(adminRouter, adminService)

	// Authentication routes
	authRouter := chi.NewRouter()
	mainRouter.Route("/auth", func(r chi.Router) {
		r.Mount("/", authRouter)
	})

	adminAuthRepository := repository.NewPostgresAdminAuthRepository(db.GetDB(), cfg.JWT)
	adminAuthService := service.NewAdminAuthService(adminAuthRepository)
	routers.SetupAuthRoutes(authRouter, adminAuthService)

	// User routes
	userRouter := chi.NewRouter()
	userRouter.Use(authMiddlewareForAdmin)
	mainRouter.Route("/api/user", func(r chi.Router) {
		r.Mount("/", userRouter)
	})

	userRepository := repository.NewPostgresUserRepository(db.GetDB())
	userService := service.NewUserService(userRepository)
	routers.SetupUserRoutes(userRouter, userService) // Set up user routes

	// Handling graceful shutdown
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
