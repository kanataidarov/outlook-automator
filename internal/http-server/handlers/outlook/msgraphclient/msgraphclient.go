package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	clientconfig "outlook-automator/internal/clientconfig"
	response "outlook-automator/pkg/api/response"
	"strconv"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	msgusers "github.com/microsoftgraph/msgraph-sdk-go/users"
)

type Response struct {
	response.Response
	Data string `json:"data,omitempty"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.outlook.msgraphclient.New"
		clCfg := clientconfig.Load()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		folders(w, r, clCfg, log)
	}
}

func folders(w http.ResponseWriter, r *http.Request, cfg *clientconfig.ClientConfig, log *slog.Logger) {
	// appGraphId := cfg.OutlookClient.AppGraphId
	clientId := cfg.OutlookClient.ClientId
	tenantId := cfg.OutlookClient.TenantId
	uname := cfg.OutlookClient.Uname
	paswd := cfg.OutlookClient.Paswd

	cred, err := azidentity.NewUsernamePasswordCredential(
		tenantId,
		clientId,
		uname,
		paswd,
		nil,
	)
	if err != nil {
		responseError(w, r, fmt.Sprintf("Error creating credentials: %v", err))
		return
	}

	gClient, err := msgraph.NewGraphServiceClientWithCredentials(cred, []string{"Mail.ReadWrite"})
	if err != nil {
		responseError(w, r, fmt.Sprintf("Error creating client: %v", err))
		return
	}

	requestFilter := ""

	requestParams := &msgusers.ItemMessagesRequestBuilderGetQueryParameters{
		Filter: &requestFilter,
		Select: []string{"subject", "sender", "receivedDateTime"},
	}

	requestConfig := &msgusers.ItemMessagesRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParams,
	}

	messages, err := gClient.Me().Messages().Get(context.Background(), requestConfig)
	if err != nil {
		log.Error("Error details", slog.String("errBody", err.Error()))
		responseError(w, r, fmt.Sprintf("Error getting messages: %v", err))
		return
	}

	responseOk(w, r, strconv.FormatInt(*messages.GetOdataCount(), 10))
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
