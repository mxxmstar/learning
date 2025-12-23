package grpc

import (
	"context"

	"github.com/mxxmstar/learning/verify_server/internal/service"
	pb "github.com/mxxmstar/learning/verify_server/proto"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthService(authService *service.AuthService) *AuthService {
	return &AuthService{
		authService: authService,
	}
}

func (s *AuthService) VerifySession(ctx context.Context, req *pb.VerifySessionRequest) (*pb.VerifySessionResponse, error) {
	user, err := s.authService.GetSessionUser(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}

}
