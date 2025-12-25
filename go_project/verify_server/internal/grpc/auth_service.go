package grpc

import (
	"context"

	"github.com/mxxmstar/learning/verify_server/internal/domain"
	"github.com/mxxmstar/learning/verify_server/internal/service"
	pb "github.com/mxxmstar/learning/verify_server/proto"
)

type AuthService struct {
	pb.UnimplementedAuthServer
	authService *service.AuthService
}

func NewAuthService(authService *service.AuthService) *AuthService {
	return &AuthService{
		authService: authService,
	}
}

func (s *AuthService) VerifySession(ctx context.Context, req *pb.VerifySessionRequest) (*pb.VerifySessionResponse, error) {
	user, err := s.authService.GetSessionUser(ctx, req.GetSessionId())
	if err != nil {
		return &pb.VerifySessionResponse{
			Valid:  false,
			UserId: 0,
			Error:  err.Error(),
		}, nil
	}

	if user == nil {
		return &pb.VerifySessionResponse{
			Valid:  false,
			UserId: 0,
			Error:  "Session not found or invalid",
		}, nil
	}

	return &pb.VerifySessionResponse{
		Valid:  true,
		UserId: user.Id,
		Error:  "",
	}, nil

}

func (s *AuthService) VerifyJWT(ctx context.Context, req *pb.VerifyJWTRequest) (*pb.VerifyJWTResponse, error) {
	claims, err := s.authService.ValidateAndParseJWT(req.GetJwtToken())
	if err != nil {
		return &pb.VerifyJWTResponse{
			Valid:    false,
			UserId:   0,
			DeviceId: "",
			Error:    err.Error(),
		}, nil
	}

	return &pb.VerifyJWTResponse{
		Valid:    true,
		UserId:   claims.UserID,
		DeviceId: claims.DeviceID,
		Error:    "",
	}, nil
}

func (s *AuthService) RefreshSession(ctx context.Context, req *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	err := s.authService.RefreshSession(ctx, req.GetSessionId())
	if err != nil {
		return &pb.RefreshSessionResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.RefreshSessionResponse{
		Success: true,
		Error:   "",
	}, nil
}

func (s *AuthService) LoginByEmail(ctx context.Context, req *pb.LoginByEmailRequest) (*pb.LoginByEmailResponse, error) {
	loginCtx := &domain.LoginContext{
		DeviceId: req.GetDeviceId(),
		// todo: 其他登录上下文信息
	}

	sessionID, err := s.authService.LoginByEmail(ctx, req.GetEmail(), req.GetPassword(), loginCtx)
	if err != nil {
		return &pb.LoginByEmailResponse{
			SessionId: "",
			JwtToken:  "",
			UserId:    0,
			Error:     err.Error(),
		}, nil
	}

	// 获取用户信息
	user, err := s.authService.GetSessionUser(ctx, sessionID)
	if err != nil {
		return &pb.LoginByEmailResponse{
			SessionId: "",
			JwtToken:  "",
			UserId:    0,
			Error:     "Failed to get user information by session ID: " + err.Error(),
		}, nil
	}

	// 生成JWT令牌
	jwtToken, err := s.authService.GenerateJWT(user, loginCtx)
	if err != nil {
		return &pb.LoginByEmailResponse{
			SessionId: "",
			JwtToken:  "",
			UserId:    0,
			Error:     "Failed to generate JWT token: " + err.Error(),
		}, nil
	}

	return &pb.LoginByEmailResponse{
		SessionId: sessionID,
		JwtToken:  jwtToken,
		UserId:    user.Id,
		Error:     "",
	}, nil
}

func (s *AuthService) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	// 验证确认密码
	if req.GetPassword() != req.GetConfirmPassword() {
		return &pb.SignUpResponse{
			Success: false,
			Error:   "Password and confirm password do not match",
		}, nil
	}

	user := &domain.User{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	err := s.authService.Signup(ctx, user)
	if err != nil {
		return &pb.SignUpResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.SignUpResponse{
		Success: true,
		Error:   "",
	}, nil
}
