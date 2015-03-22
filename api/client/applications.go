package client

import (
	"net/http"

	"github.com/kiasaki/hazel/api/data"
)

type ApplicationsClient struct {
	Client Client
}

func (c Client) Applications() ApplicationsClient {
	return ApplicationsClient{c}
}

type ApplicationsResponse struct {
	Applications []data.Application `json:"applications"`
}

func (c ApplicationsClient) All() (ApplicationsResponse, error) {
	var response ApplicationsResponse
	err := c.Client.RequestAndProcess("GET", "/applications", nil, &response, http.StatusOK)
	return response, err
}
