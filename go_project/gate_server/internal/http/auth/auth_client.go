package http_auth_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	common_auth "github.com/mxxmstar/learning/pkg/common/auth"
)

type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthClient(baseURL string, httpClient *http.Client) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *AuthClient) VerifySession(ctx context.Context, sessionID string) (*common_auth.VerifySessionResponse, error) {
	req := &common_auth.VerifySessionRequest{
		SessionID: sessionID,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Post(
		fmt.Sprintf("%s/gate/user-auth/verify-session", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res common_auth.VerifySessionResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *AuthClient) VerifyJWT(ctx context.Context, jwt string) (*common_auth.VerifyJWTResponse, error) {
	req := &common_auth.VerifyJWTRequest{
		JWTToken: jwt,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Post(
		fmt.Sprintf("%s/gate/user-auth/verify-jwt", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res common_auth.VerifyJWTResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *AuthClient) RefreshSession(ctx context.Context, sessionID string) (*common_auth.RefreshSessionResponse, error) {
	req := &common_auth.RefreshSessionRequest{
		SessionID: sessionID,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Post(
		fmt.Sprintf("%s/gate/user-auth/refresh-session", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res common_auth.RefreshSessionResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
