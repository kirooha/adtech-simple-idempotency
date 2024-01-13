package jobhandlers

import (
	"context"
	"encoding/json"

	"adtech.simple/internal/pkg/adserverclient"
	"adtech.simple/internal/pkg/dbquery"
	"adtech.simple/internal/pkg/model"
	"github.com/vgarvardt/gue/v5"
)

type CampaignCreationHandler struct {
	adServerClient *adserverclient.Client
	dbQuerier      *dbquery.Queries
}

func NewCampaignCreationHandler(adServerClient *adserverclient.Client, dbQuerier *dbquery.Queries) *CampaignCreationHandler {
	return &CampaignCreationHandler{
		adServerClient: adServerClient,
		dbQuerier:      dbQuerier,
	}
}

func (h *CampaignCreationHandler) MakeHandler() func(ctx context.Context, j *gue.Job) error {
	return func(ctx context.Context, j *gue.Job) error {
		var params model.JobParamCampaignCreationEvent
		if err := json.Unmarshal(j.Args, &params); err != nil {
			return err
		}

		campaign, err := h.dbQuerier.GetCampaign(ctx, params.ID)
		if err != nil {
			return err
		}

		adServerID, err := h.adServerClient.CreateCampaign(ctx, campaign.Name, campaign.Description)
		if err != nil {
			return err
		}

		updatingParams := dbquery.UpdateCampaignParams{ID: params.ID, AdserverID: adServerID}
		_, err = h.dbQuerier.UpdateCampaign(ctx, updatingParams)
		if err != nil {
			return err
		}

		return nil
	}
}
