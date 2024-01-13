package adserverclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	createCampaignURL = "/create-campaign"
)

type campaignCreationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type campaignCreationResponse struct {
	ID uuid.UUID
}

func (c *Client) CreateCampaign(ctx context.Context, name, description string) (*uuid.UUID, error) {
	var (
		req  campaignCreationRequest
		resp campaignCreationResponse
	)

	req.Name = name
	req.Description = description

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", baseURL, createCampaignURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is %d", httpResp.StatusCode)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp.ID, nil
}
