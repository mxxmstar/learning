package http_auth_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mxxmstar/learning/gate_server/gate_config"
	http_status_client "github.com/mxxmstar/learning/gate_server/internal/http/status"
	auth_def "github.com/mxxmstar/learning/pkg/def/auth"
)

type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthClient(c gate_config.Config, httpClient *http.Client) (*AuthClient, error) {
	sc := http_status_client.NewStatusClient(c, httpClient)
	verifyServer, err := http_status_client.GetVerifyServer(sc)
	if err != nil {
		return nil, err
	}
	return &AuthClient{
		baseURL:    verifyServer.HTTPAddress.Host + ":" + fmt.Sprint(verifyServer.HTTPAddress.Port),
		httpClient: httpClient,
	}, nil
}

func (c *AuthClient) VerifySession(ctx context.Context, sessionId string) (*auth_def.VerifySessionResponse, error) {
	req := &auth_def.VerifySessionRequest{
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

	var res auth_def.VerifySessionResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *AuthClient) VerifyJWT(ctx context.Context, jwt string) (*auth_def.VerifyJWTResponse, error) {
	req := &auth_def.VerifyJWTRequest{
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

	var res auth_def.VerifyJWTResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *AuthClient) RefreshSession(ctx context.Context, sessionId string) (*auth_def.RefreshSessionResponse, error) {
	req := &auth_def.RefreshSessionRequest{
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

	var res auth_def.RefreshSessionResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *AuthClient) LoginByEmail(ctx context.Context, email, password, DeviceId string) (*auth_def.LoginByEmailResponse, error) {
	req := &auth_def.LoginByEmailRequest{
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

	var res auth_def.LoginByEmailResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *AuthClient) SignUp(ctx context.Context, username, email, password, confirm_password string) (*auth_def.SignUpResponse, error) {
	req := &auth_def.SignUpRequest{
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

	var res auth_def.SignUpResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
