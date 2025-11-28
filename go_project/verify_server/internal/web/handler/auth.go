package handler

import (
	"log"
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/mxxmstar/learning/pkg/logger"
	"github.com/mxxmstar/learning/verify_server/internal/domain"
	"github.com/mxxmstar/learning/verify_server/internal/service"
)

type AuthAndler struct {
	authService *service.AuthService
	userService *service.UserService
	// emailExp 邮箱正则表达式
	emailExp *regexp.Regexp
	// passwordExp 密码正则表达式
	passwordExp *regexp.Regexp
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthAndler {
	return &AuthAndler{
		authService: authService,
		userService: userService,
		emailExp:    regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, regexp.None),
		passwordExp: regexp.MustCompile(`^[a-zA-Z0-9_-]{6,20}$`, regexp.None),
	}
}

func (h *AuthAndler) SignupHandler(ctx *gin.Context) {
	// println("SignupHandler")
	type SignupRequest struct {
		Email           string `json:"email"`
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	var req SignupRequest
	// Bind 方法会根据 Content-Type 的不同，使用不同的绑定方法
	// 解析错误会直接返回 400 错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := h.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
	}

	if !ok {
		ctx.String(http.StatusOK, "email format error")
		return
	}

	ok, err = h.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
	}

	if !ok {
		ctx.String(http.StatusOK, "password format error")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "password confirmation does not match")
		return
	}

	// 数据库操作
	err = h.authService.Signup(ctx, &domain.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err == service.ErrUserEmailConflict {
		ctx.String(http.StatusOK, "email already has been registered.")
		return
	}
	if err == service.ErrUserUsernameConflict {
		ctx.String(http.StatusOK, "username already has been registered.")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "system error.")
		return
	}

	ctx.String(http.StatusOK, "signup success.")
	log.Println("signup success")
	logger.LogAuth(ctx, "signup", true, "signup success")
}

func (h *AuthAndler) LoginHandler(ctx *gin.Context) {
	println("LoginHandler")
}

func (h *AuthAndler) OAuthHandler(ctx *gin.Context) {
	println("OAuthHandler")
}
