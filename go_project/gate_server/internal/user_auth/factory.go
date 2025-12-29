package auth_user

import (
	auth_client "github.com/mxxmstar/learning/gate_server/internal/grpc/auth"
	http_auth_client "github.com/mxxmstar/learning/gate_server/internal/http/auth"
)

type AuthUserType string

const (
	AuthUserGRPC AuthUserType = "grpc"
	AuthUserHTTP AuthUserType = "http"
)

func NewAuthService(authUserType AuthUserType, g *auth_client.AuthClient, h *http_auth_client.AuthClient) AuthService {
	switch authUserType {
	case AuthUserGRPC:
		return NewGRPCAuthService(g)
	case AuthUserHTTP:
		return NewHTTPAuthService(h)
	default:
		return NewGRPCAuthService(g)
	}
}
