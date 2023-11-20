package main

import (
	"log/slog"
	"os"
	"user-admin/internal/config"
	"user-admin/pkg/database"
	log_utils "user-admin/pkg/lib/logger_utils"
	"user-admin/pkg/logger"
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
}
