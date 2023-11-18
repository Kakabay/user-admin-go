package main

import (
	"log/slog"
	"user-admin/internal/config"
	"user-admin/internal/logger"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.SetupLogger(cfg.Env)

	log.Info("Starting the server...", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled") // if env is set to prod, debug messages are going to be disabled
}
