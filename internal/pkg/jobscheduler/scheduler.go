package jobscheduler

import (
	"context"
	"encoding/json"
	"fmt"

	"adtech.simple/internal/pkg/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vgarvardt/gue/v5"
	"github.com/vgarvardt/gue/v5/adapter/pgxv5"
)

type Scheduler struct {
	gueClient *gue.Client
}

func NewScheduler(gueClient *gue.Client) *Scheduler {
	return &Scheduler{
		gueClient: gueClient,
	}
}

func (s *Scheduler) ScheduleCreateCampaignInAdServerTx(ctx context.Context, tx pgx.Tx, id *uuid.UUID) error {
	jobParamCampaignCreationEvent := model.JobParamCampaignCreationEvent{
		ID: id,
	}
	b, err := json.Marshal(jobParamCampaignCreationEvent)
	if err != nil {
		return err
	}

	j := &gue.Job{
		Type:  fmt.Sprint(model.JobTypeCampaignCreation),
		Queue: fmt.Sprint(model.QueueTypeCampaignCreation),
		Args:  b,
	}
	if err := s.gueClient.EnqueueTx(ctx, j, pgxv5.NewTx(tx)); err != nil {
		return err
	}

	return nil
}
