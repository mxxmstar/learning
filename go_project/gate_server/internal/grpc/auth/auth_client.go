package grpc_auth_client

import (
	"context"
	"fmt"

	"github.com/mxxmstar/learning/gate_server/gate_config"
	"github.com/mxxmstar/learning/pkg/logger"
	pb "github.com/mxxmstar/learning/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client pb.AuthClient
	config *gate_config.Config
}

func NewAuthClient(config *gate_config.Config) (*AuthClient, error) {
	// 连接到 verify_server 的 gRPC 服务
	dsn := fmt.Sprintf("%s:%d", config.VerifyServer.VerifyGRPCServerConfig.Host, config.VerifyServer.VerifyGRPCServerConfig.Port)
	conn, err := grpc.NewClient(
		dsn,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.FormatLog(context.Background(), "error", fmt.Sprintf("Failed to connect to verify_server gRPC server at %s: %v", dsn, err))
		return nil, err
	}

	client := pb.NewAuthClient(conn)

	authClient := &AuthClient{
		conn:   conn,
		client: client,
		config: config,
	}
	return authClient, nil
}

// 验证 session
func (c *AuthClient) VerifySession(ctx context.Context, sessionId string) (*pb.VerifySessionResponse, error) {
	req := &pb.VerifySessionRequest{
		SessionId: sessionId,
	}
	return c.client.VerifySession(ctx, req)
}

// 验证 JWT
func (c *AuthClient) VerifyJWT(ctx context.Context, jwt string) (*pb.VerifyJWTResponse, error) {
	req := &pb.VerifyJWTRequest{
		JwtToken: jwt,
	}
	return c.client.VerifyJWT(ctx, req)
}

// 刷新 session
func (c *AuthClient) RefreshSession(ctx context.Context, sessionId string) (*pb.RefreshSessionResponse, error) {
	req := &pb.RefreshSessionRequest{
		SessionId: sessionId,
	}
	return c.client.RefreshSession(ctx, req)
}

// 通过邮箱登录
func (c *AuthClient) LoginByEmail(ctx context.Context, email, password, DeviceId string) (*pb.LoginByEmailResponse, error) {
	req := &pb.LoginByEmailRequest{
		Email:    email,
		Password: password,
		DeviceId: DeviceId,
	}
	return c.client.LoginByEmail(ctx, req)
}

// 注册
func (c *AuthClient) SignUp(ctx context.Context, username, email, password, confirmPassword string) (*pb.SignUpResponse, error) {
	req := &pb.SignUpRequest{
		Username:        username,
		Email:           email,
		Password:        password,
		ConfirmPassword: confirmPassword,
	}
	return c.client.SignUp(ctx, req)
}

// 关闭客户端连接
func (c *AuthClient) Close() error {
	return c.conn.Close()
}
