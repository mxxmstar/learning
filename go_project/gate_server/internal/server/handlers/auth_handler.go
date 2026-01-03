package handlers

import (
	"github.com/mxxmstar/learning/gate_server/gate_config"
	auth_user "github.com/mxxmstar/learning/gate_server/internal/user_auth"
)

type AuthHandler struct {
	authService auth_user.AuthService
}

func NewAuthHandler(cfg *gate_config.Config) *AuthHandler {
	var authService auth_user.AuthService

}
