package handlers

import (
	auth_user "github.com/mxxmstar/learning/gate_server/internal/user_auth"
)

type AuthMessageHandler struct {
	authService auth_user.AuthService
}

func NewAuthMessageHandler(authService auth_user.AuthService) *AuthMessageHandler {
	return &AuthMessageHandler{
		authService: authService,
	}
}
