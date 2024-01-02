package main

import (
	"log/slog"
	"os"
	"outlook-automator/internal/http-server/router"
	"outlook-automator/pkg/config"
)

func main() {
	cfg := config.Load()
	log := setLogger(cfg.LogLevel)

	router.Serve(cfg, log)
}

func setLogger(logLevel string) *slog.Logger {
	var logger *slog.Logger

	switch logLevel {
	case "debug":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "info":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}

	return logger
}
