package main

import (
	"log/slog"
	"os"
	"outlook-automator/internal/http-server/router"
	"outlook-automator/pkg/server/config"
)

func main() {
	srvCfg := config.Load()
	log := setLogger(srvCfg.LogLevel)

	router.Serve(srvCfg, log)
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
