package adserverclient

import (
	"net/http"
)

const baseURL = "http://localhost:9091"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: http.DefaultClient,
	}
}
