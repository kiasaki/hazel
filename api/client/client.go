package client

import (
	"net/http"

	"github.com/kiasaki/batbelt/rest"
)

type Client struct {
	ApiBaseUrl string
	JwtToken   string
}

func NewClient(baseUrl, token string) Client {
	return Client{baseUrl, token}
}

func (c Client) Request(method, url string, entity interface{}) (*http.Response, error) {
	return rest.MakeRequestWithMiddleware(method, c.ApiBaseUrl+url, entity, func(r *http.Request) {
		if c.JwtToken != "" {
			r.Header.Set("Authorization", "Bearer "+c.JwtToken)
		}
	})
}

func (c Client) RequestAndProcess(method, url string, entity interface{}, responseEntity interface{}, expectedStatus int) error {
	response, err := c.Request(method, url, entity)
	if err != nil {
		return err
	}
	return rest.ProcessResponseEntity(response, responseEntity, expectedStatus)
}
