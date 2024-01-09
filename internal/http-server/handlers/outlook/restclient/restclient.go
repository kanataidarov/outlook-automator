package restclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	clientconfig "outlook-automator/internal/clientconfig"
	response "outlook-automator/pkg/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Data string `json:"data,omitempty"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.outlook.restclient.New"
		clCfg := clientconfig.Load()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		folders(w, r, clCfg, log)
	}
}

func folders(w http.ResponseWriter, r *http.Request, clCfg *clientconfig.ClientConfig, log *slog.Logger) {
	selectField := r.URL.Query().Get("field")
	token := clCfg.OutlookClient.Paswd // this config field used to store token in case of the RestClient
	requestUrl := fmt.Sprintf("https://graph.microsoft.com/v1.0/me?$select=%s", selectField)

	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		responseError(w, r, fmt.Sprintf("Error making request: %v", err))
		return
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		responseError(w, r, fmt.Sprintf("Error retrieving response body: %v", err))
	}

	respBody := make(map[string]interface{})
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		responseError(w, r, fmt.Sprintf("Error unmarshalling response body: %v", err))
	}

	respField, _ := json.Marshal(map[string]string{
		selectField: fmt.Sprintf("%v", respBody[selectField]),
	})

	responseOk(w, r, string(respField))
}

func responseOk(w http.ResponseWriter, r *http.Request, data string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Data:     data,
	})
}

func responseError(w http.ResponseWriter, r *http.Request, errMsg string) {
	render.JSON(w, r, response.Error(errMsg))
}
