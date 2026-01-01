package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mxxmstar/learning/gate_server/gate_config"
	"github.com/mxxmstar/learning/gate_server/internal/conn"
	grpc_auth_client "github.com/mxxmstar/learning/gate_server/internal/grpc/auth"
	http_auth_client "github.com/mxxmstar/learning/gate_server/internal/http/auth"
	auth_user "github.com/mxxmstar/learning/gate_server/internal/user_auth"
)

type AuthMessageHandler struct {
	authService auth_user.AuthService
}

func NewAuthMessageHandler(authService auth_user.AuthService) *AuthMessageHandler {
	return &AuthMessageHandler{authService: authService}
}

func (h *AuthMessageHandler) HandleMessage(ctx context.Context, conn conn.Connection, envelope *Envelope) error {
	switch envelope.Type {
	case "signup":
		return h.handleSignup(ctx, conn, envelope)
	case "login":
		// return h.handleLogin(ctx, conn, envelope)
		return fmt.Errorf("login")
	default:
		return fmt.Errorf("unknown message type: %s", envelope.Type)
	}
}

// 处理用户注册，调用 grpc SignUp 方法
func (h *AuthMessageHandler) handleSignup(ctx context.Context, conn conn.Connection, envelope *Envelope) error {
	// 从 envelope 中获取用户注册信息
	email, ok := envelope.Body["email"].(string)
	if !ok {
		return fmt.Errorf("signup: email is missing or invalid")
	}

	username, ok := envelope.Body["username"].(string)
	if !ok {
		return fmt.Errorf("signup: username is missing or invalid")
	}

	password, ok := envelope.Body["password"].(string)
	if !ok {
		return fmt.Errorf("signup: password is missing or invalid")
	}

	confirmPassword, ok := envelope.Body["confirm_password"].(string)
	if !ok {
		return fmt.Errorf("signup: confirm_password is missing or invalid")
	}

	if password != confirmPassword {
		return fmt.Errorf("signup: password and confirm_password do not match")
	}

	authGRPC, ok := h.authService.(*auth_user.GRPCAuthService)
	if !ok {
		return fmt.Errorf("signup: authService is not a GRPCAuthService")
	}

	signupRsp, err := authGRPC.AuthService().SignUp(ctx, username, email, password, confirmPassword)
	if err != nil {
		return fmt.Errorf("signup: %v", err)
	}

	response := map[string]interface{}{
		"type":    "signup_response",
		"success": signupRsp.Error == "",
		"error":   signupRsp.Error,
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("signup: %v", err)
	}
	return conn.Send(responseBytes)
}

func InitWebSocketServer(cfg *gate_config.Config) *WebsocketServer {
	// 初始化 gRPC 客户端
	grpcClient, err := grpc_auth_client.NewAuthClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}

	// 初始化 HTTP 客户端
	httpClient := http_auth_client.NewAuthClient(cfg.VerifyServer.Host+fmt.Sprintf(":%d", cfg.VerifyServer.Port), &http.Client{})

	// 创建认证服务
	authService := auth_user.NewAuthService(auth_user.AuthUserGRPC, grpcClient, httpClient)

	// 创建连接管理器
	connManager := conn.NewManager()

	// 创建 WebSocket 服务器
	wsServer := NewWebsocketServer(
		cfg.GateServer.Name, // gate Id
		connManager,
		authService,
		nil, // store (可以是 Redis 或其他存储)
		nil, // notifyOld function
	)

	// 注册认证消息处理器
	authHandler := NewAuthMessageHandler(authService)
	wsServer.RegisterHandler("login", authHandler)
	wsServer.RegisterHandler("signup", authHandler)

	return wsServer
}
