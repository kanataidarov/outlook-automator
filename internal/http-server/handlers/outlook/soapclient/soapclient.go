package soapclient

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	clientconfig "outlook-automator/internal/clientconfig"
	response "outlook-automator/pkg/api/response"

	"github.com/azure/go-ntlmssp"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Data string `json:"data,omitempty"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.outlook.soapclient.New"
		clCfg := clientconfig.Load()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		folders(w, r, clCfg, log)
	}
}

func folders(w http.ResponseWriter, r *http.Request, clCfg *clientconfig.ClientConfig, log *slog.Logger) {
	uname := clCfg.OutlookClient.Uname
	paswd := clCfg.OutlookClient.Paswd
	ewsUrl := "https://owa.beeline.kz/ews/exchange.asmx"

	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}

	body := `<?xml version="1.0" encoding="utf-8"?>
	<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
	  <soap:Body>
		<FindItem xmlns="http://schemas.microsoft.com/exchange/services/2006/messages" xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types" Traversal="Shallow">
		  <ItemShape>
			<t:BaseShape>IdOnly</t:BaseShape>
		  </ItemShape>
		  <ParentFolderIds>
			<t:DistinguishedFolderId Id="deleteditems"/>
		  </ParentFolderIds>
		</FindItem>
	  </soap:Body>
	</soap:Envelope>`

	req, err := http.NewRequest("POST", ewsUrl, bytes.NewReader([]byte(body)))
	if err != nil {
		responseError(w, r, fmt.Sprintf("Error creating request: %v", err))
		return
	}
	req.SetBasicAuth(uname, paswd)
	req.Header.Set("Content-Type", "text/xml")
	defer req.Body.Close()

	resp, err := client.Do(req)
	if err != nil {
		responseError(w, r, fmt.Sprintf("Error making request: %v", err))
		return
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		responseError(w, r, fmt.Sprintf("Error parsing response: %v", err))
		return
	}

	responseOk(w, r, string(respBytes[:]))
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
