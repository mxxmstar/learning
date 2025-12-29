package auth_user

import (
	"context"

	grpc_auth_client "github.com/mxxmstar/learning/gate_server/internal/grpc/auth"
	http_auth_client "github.com/mxxmstar/learning/gate_server/internal/http/auth"
)

// 验证服务接口
type AuthService interface {
	ValidateTokenOrSession(ctx context.Context, token, sessionID, deviceID string) (*AuthResult, error)
	RefreshSession(ctx context.Context, sessionID string) (*AuthResult, error)
}

// AuthResult 认证结果
type AuthResult struct {
	UserID   uint64
	DeviceID string
	Valid    bool
	Error    string
}

type GRPCAuthService struct {
	authService *grpc_auth_client.AuthClient
}

type HTTPAuthService struct {
	authService *http_auth_client.AuthClient
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
