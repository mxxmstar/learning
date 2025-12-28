package auth_user

import (
	grpc_client "github.com/mxxmstar/learning/gate_server/internal/grpc"
	http_client "github.com/mxxmstar/learning/gate_server/internal/http"
)

type AuthUserType string

const (
	AuthUserGRPC AuthUserType = "grpc"
	AuthUserHTTP AuthUserType = "http"
)

func NewAuthService(authUserType AuthUserType, g *grpc_client.AuthClient, h *http_client.AuthClient) AuthService {
	switch authUserType {
	case AuthUserGRPC:
		return NewGRPCAuthService(g)
	case AuthUserHTTP:
		return NewHTTPAuthService(h)
	default:
		return NewGRPCAuthService(g)
	}
}
