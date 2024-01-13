package model

import "github.com/google/uuid"

//go:generate stringer -type=JobType
type JobType int16

const (
	JobTypeCampaignCreation JobType = 1
)

//go:generate stringer -type=QueueType
type QueueType int16

const (
	QueueTypeCampaignCreation QueueType = 1
)

type JobParamCampaignCreationEvent struct {
	ID *uuid.UUID
}
