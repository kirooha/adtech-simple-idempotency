package adtechsimpleapi

import (
	"encoding/json"
	"log"
	"net/http"

	"adtech.simple/internal/pkg/adserverclient"
	"github.com/google/uuid"
)

type AddViewHandler struct {
	adServerClient *adserverclient.Client
}

type addViewRequest struct {
	EventID uuid.UUID `json:"event_id"`
	AdID    uuid.UUID `json:"ad_id"`
}

func NewAddViewHandler(adServerClient *adserverclient.Client) *AddViewHandler {
	return &AddViewHandler{
		adServerClient: adServerClient,
	}
}

func (h *AddViewHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var (
		msgPrefix = "adtechsimpleapi.AddViewHandler.ServeHTTP"
		ctx       = request.Context()
		req       addViewRequest
	)
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		log.Printf("%s: json.NewDecoder(request.Body).Decode error: %v", msgPrefix, err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.adServerClient.AddView(ctx, req.EventID, req.AdID); err != nil {
		log.Printf("%s: h.adServerClient.AddView error: %v", msgPrefix, err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}
