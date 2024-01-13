package adtechsimpleapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"adtech.simple/internal/pkg/adserverclient"
	"adtech.simple/internal/pkg/dbquery"
	"adtech.simple/internal/pkg/jobscheduler"
	"adtech.simple/internal/pkg/store"
	"github.com/jackc/pgx/v5"
)

type CreateCampaignHandler struct {
	storage        *store.Storage
	adServerClient *adserverclient.Client
	jobScheduler   *jobscheduler.Scheduler
}

type createCampaignHandlerRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewCreateCampaignHandler(storage *store.Storage, adServerClient *adserverclient.Client, jobScheduler *jobscheduler.Scheduler) *CreateCampaignHandler {
	return &CreateCampaignHandler{
		storage:        storage,
		adServerClient: adServerClient,
		jobScheduler:   jobScheduler,
	}
}

func (h *CreateCampaignHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var (
		msgPrefix = "adtechsimpleapi.CreateCampaignHandler.ServeHTTP"
		ctx       = request.Context()
		req       createCampaignHandlerRequest
	)
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		log.Printf("%s: json.NewDecoder(request.Body).Decode error: %v", msgPrefix, err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := h.storage.WithTx(ctx, func(ctx context.Context, tx pgx.Tx, queries *dbquery.Queries) error {
		createCampaignParams := dbquery.CreateCampaignParams{
			Name:        req.Name,
			Description: req.Description,
		}
		createdCampaign, err := queries.CreateCampaign(ctx, createCampaignParams)
		if err != nil {
			return err
		}

		return h.jobScheduler.ScheduleCreateCampaignInAdServerTx(ctx, tx, createdCampaign.ID)
	})
	if err != nil {
		log.Printf("%s: h.storage.WithTx error: %v", msgPrefix, err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}
