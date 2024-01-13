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
	addViewURL = "/add-view"
)

type addViewRequest struct {
	EventID uuid.UUID `json:"event_id"`
	AdID    uuid.UUID `json:"ad_id"`
}

func (c *Client) AddView(ctx context.Context, eventID, adID uuid.UUID) error {
	var (
		req addViewRequest
	)

	req.EventID = eventID
	req.AdID = adID

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s", baseURL, addViewURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is %d", httpResp.StatusCode)
	}

	return nil
}
