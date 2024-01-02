package router

import (
	"log/slog"
	"net/http"
	"outlook-automator/internal/http-server/handlers/outlook/restclient"
	"outlook-automator/pkg/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Serve(cfg *config.Config, log *slog.Logger) {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/folders", restclient.New(cfg, log))

	log.Info("Starting server", slog.String("address", cfg.Address))
	log.Debug("Debug logs enabled")

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server", slog.String("address", cfg.Address))
	}

	log.Error("Server stopped")
}
