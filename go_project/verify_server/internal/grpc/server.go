package grpc_server

import (
	"context"
	"fmt"
	"net"

	"github.com/mxxmstar/learning/pkg/logger"
	pb "github.com/mxxmstar/learning/proto"
	"github.com/mxxmstar/learning/verify_server/internal/service"
	"github.com/mxxmstar/learning/verify_server/verify_config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCService struct {
	authService  *service.AuthService
	useerService *service.UserService
}

// GRPCServer gRPC 服务器结构体
type GRPCServer struct {
	grpcService *GRPCService
	config      *verify_config.Config
	server      *grpc.Server
}

func NewGRPCServer(grpcService *GRPCService, config *verify_config.Config) *GRPCServer {
	return &GRPCServer{
		grpcService: grpcService,
		config:      config,
	}
}

func (s *GRPCServer) Start() error {
	// 创建监听地址
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.VerifyGRPCServer.Port))
	if err != nil {
		logger.FormatLog(context.Background(), "error", fmt.Sprintf("Failed to start gRPC server at %d: %v", s.config.VerifyGRPCServer.Port, zap.Error(err)))
		return err
	}

	// 创建 gRPC 服务器
	s.server = grpc.NewServer()

	// 注册服务
	authService := NewAuthService(s.grpcService.authService)
	// userService := NewUserService(s.useerService)
	pb.RegisterAuthServer(s.server, authService)

	// 在开发环境中启用反射服务，以便使用 gRPC 客户端工具进行调试
	if s.config.ProjectConfig.Env != "production" {
		reflection.Register(s.server)
		logger.FormatLog(context.Background(), "info", "gRPC reflection service enabled in non-production environment")
	}

	// 启动服务器
	if err := s.server.Serve(listen); err != nil {
		logger.FormatLog(context.Background(), "error", fmt.Sprintf("Failed to start gRPC server: %v", err))
		return err
	}
	return nil
}

func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
		logger.FormatLog(context.Background(), "info", "gRPC server stopped gracefully")
	}
}
