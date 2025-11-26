package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mxxmstar/learning/verify_server/internal/service"
)

type AuthAndler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthAndler {
	return &AuthAndler{
		authService: authService,
		userService: userService,
	}
}

func (h *AuthAndler) SignupHandler(ctx *gin.Context) {

}

func (h *AuthAndler) LoginHandler(ctx *gin.Context) {

}

func (h *AuthAndler) OAuthHandler(ctx *gin.Context) {

}
