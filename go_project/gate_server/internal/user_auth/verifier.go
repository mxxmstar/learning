package auth_user

import (
	"context"

	grpc_auth_client "github.com/mxxmstar/learning/gate_server/internal/grpc/auth"
	http_auth_client "github.com/mxxmstar/learning/gate_server/internal/http/auth"
)

// 验证服务接口
type AuthService interface {
	ValidateTokenOrSession(ctx context.Context, token, sessionId, deviceId string) (*AuthResult, error)
	RefreshSession(ctx context.Context, sessionId string) (*AuthResult, error)
}

// AuthResult 认证结果
type AuthResult struct {
	UserId   uint64
	DeviceId string
	Valid    bool
	Error    string
}

type GRPCAuthService struct {
	authService *grpc_auth_client.AuthClient
}

func (g *GRPCAuthService) AuthService() *grpc_auth_client.AuthClient {
	return g.authService
}

type HTTPAuthService struct {
	authService *http_auth_client.AuthClient
}

func (h *HTTPAuthService) AuthService() *http_auth_client.AuthClient {
	return h.authService
}

func NewGRPCAuthService(authService *grpc_auth_client.AuthClient) *GRPCAuthService {
	return &GRPCAuthService{
		authService: authService,
	}
}

func NewHTTPAuthService(authService *http_auth_client.AuthClient) *HTTPAuthService {
	return &HTTPAuthService{
		authService: authService,
	}
}
