package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"user-admin/internal/config"
	handlers "user-admin/internal/delivery/v1/handlers/user"
	repository "user-admin/internal/repository/postgres"
	"user-admin/internal/service"
	database "user-admin/pkg/database"
	log_utils "user-admin/pkg/lib/logger_utils"
	"user-admin/pkg/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.SetupLogger(cfg.Env)

	log.Info("Starting the server...", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled") // If env is set to prod, debug messages are going to be disabled

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Error("Failed to init database:", log_utils.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	userRepository := repository.NewPostgresUserRepository(db.GetDB())
	userService := service.NewUserService(userRepository)

	userHandler := handlers.NewUserHandler(userService)

	// Handle graceful shutdown by using concurent function that is going to wait signal in backround
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Info("Shutting down the server gracefully...")

		if err := db.Close(); err != nil {
			log.Error("Error closing database:", log_utils.Err(err))
		}
		os.Exit(0)
	}()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Mount routes from userHandler
	router.Mount("/api", userHandler.NewRoutes())

	err = http.ListenAndServe(":8082", router)
	if err != nil {
		log.Error("Server failed to start:", err)
	}
}
