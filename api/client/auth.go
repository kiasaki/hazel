package client

import (
	"net/http"
)

type AuthClient struct {
	Client Client
}

func (c Client) Auth() AuthClient {
	return AuthClient{c}
}

// Attempts login againt Hazel API and returns JWT token on success
func (c AuthClient) Login(email, password string) (string, error) {
	var loginRequest = struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{email, password}
	var loginResponse struct {
		Token string `json:"token"`
	}
	err := c.Client.RequestAndProcess("POST", "/auth/login", loginRequest, &loginResponse, http.StatusOK)
	return loginResponse.Token, err
}
