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

func (c *AuthClient) VerifySession(ctx context.Context, sessionId string) (*common_auth.VerifySessionResponse, error) {
	req := &common_auth.VerifySessionRequest{
		SessionId: sessionId,
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

func (c *AuthClient) RefreshSession(ctx context.Context, sessionId string) (*common_auth.RefreshSessionResponse, error) {
	req := &common_auth.RefreshSessionRequest{
		SessionId: sessionId,
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

func (c *AuthClient) LoginByEmail(ctx context.Context, email, password, DeviceId string) (*common_auth.LoginByEmailResponse, error) {
	req := &common_auth.LoginByEmailRequest{
		Email:    email,
		Password: password,
		DeviceId: DeviceId,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Post(
		fmt.Sprintf("%s/gate/user-auth/loginByEmail", c.baseURL),
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

	var res common_auth.LoginByEmailResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *AuthClient) SignUp(ctx context.Context, username, email, password, confirm_password string) (*common_auth.SignUpResponse, error) {
	req := &common_auth.SignUpRequest{
		Username:        username,
		Email:           email,
		Password:        password,
		ConfirmPassword: confirm_password,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Post(
		fmt.Sprintf("%s/gate/user-auth/signup", c.baseURL),
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

	var res common_auth.SignUpResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
