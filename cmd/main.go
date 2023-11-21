package main

import (
	"log/slog"
	"net/http"
	"os"
	"user-admin/internal/config"
	"user-admin/pkg/database"
	log_utils "user-admin/pkg/lib/logger_utils"
	"user-admin/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.SetupLogger(cfg.Env)

	log.Info("Starting the server...", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled") // if env is set to prod, debug messages are going to be disabled

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Error("Failed to init database:", log_utils.Err(err))
		os.Exit(1)
	} 
	defer db.Close()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	
	log.Info("Starting the server...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Error("Server failed to start:", err)
	}
}
