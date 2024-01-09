package router

import (
	"log/slog"
	"net/http"
	client "outlook-automator/internal/http-server/handlers/outlook/soapclient"
	"outlook-automator/pkg/server/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Serve(srvCfg *config.ServerConfig, log *slog.Logger) {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/folders", client.New(log))

	log.Info("Starting server", slog.String("address", srvCfg.Address))
	log.Debug("Debug logs enabled")

	srv := &http.Server{
		Addr:         srvCfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  srvCfg.HttpServer.Timeout,
		WriteTimeout: srvCfg.HttpServer.Timeout,
		IdleTimeout:  srvCfg.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server", slog.String("address", srvCfg.HttpServer.Address))
	}

	log.Error("Server stopped")
}
