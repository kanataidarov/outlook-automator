package client

import (
	"log/slog"
	"net/http"
	response "outlook-automator/pkg/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.outlook.client.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID((r.Context()))),
		)

		responseOk(w, r)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: response.OK(),
	})
}

func responseError(w http.ResponseWriter, r *http.Request, errMsg string) {
	render.JSON(w, r, response.Error(errMsg))
}
