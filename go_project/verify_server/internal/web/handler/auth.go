package handler

import (
	"log"
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/mxxmstar/learning/pkg/logger"
	"github.com/mxxmstar/learning/verify_server/internal/domain"
	"github.com/mxxmstar/learning/verify_server/internal/service"
	"github.com/mxxmstar/learning/verify_server/internal/web/response"
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
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignupRequest
	// Bind 方法会根据 Content-Type 的不同，使用不同的绑定方法
	// 解析错误会直接返回 400 错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := h.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, response.ErrorResponse("system error", nil))
		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, response.ErrorResponse("email format error", nil))
		return
	}

	ok, err = h.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, response.ErrorResponse("system error", nil))
		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, response.ErrorResponse("password format error", nil))
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusOK, response.ErrorResponse("password confirmation does not match", nil))
		return
	}

	// 数据库操作
	err = h.authService.Signup(ctx, &domain.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err == service.ErrUserEmailConflict {
		ctx.JSON(http.StatusOK, response.ErrorResponse("email already has been registered.", nil))
		return
	}
	if err == service.ErrUserUsernameConflict {
		ctx.JSON(http.StatusOK, response.ErrorResponse("username already has been registered.", nil))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusOK, response.ErrorResponse("system error.", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse("signup success.", nil))
	log.Println("signup success")
	logger.LogAuth(ctx, "signup", true, "signup success")
}

func (h *AuthAndler) LoginHandler(ctx *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		DeviceId string `json:"deviceId"`
	}

	var req LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 创建登录上下文
	loginCtx := &domain.LoginContext{
		DeviceId:  req.DeviceId,
		IPAddress: ctx.ClientIP(),
		// UserAgent: ctx.GetHeader("User-Agent"),
		UserAgent: ctx.Request.UserAgent(),
	}

	// 传统 session 登录方式
	sessionID, err := h.authService.LoginByEmail(ctx, req.Email, req.Password, loginCtx)
	if err != nil {
		ctx.JSON(http.StatusOK, response.ErrorResponse("invalid username or password", nil))
		return
	}

	// 获取用户信息生成 JWT 令牌
	user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, response.ErrorResponse("login success but failed to generate jwt token", nil))
		return
	}

	jwtToken, err := h.authService.GenerateJWT(user, loginCtx)
	if err != nil {
		ctx.JSON(http.StatusOK, response.ErrorResponse("login success but failed to generate jwt token", nil))
		return
	}

	responseData := map[string]interface{}{
		"sessionId": sessionID,
		"jwtToken":  jwtToken,
		"userId":    user.Id,
	}

	ctx.Header("x-jwt-token", jwtToken) // 将 JWT 令牌添加到响应头中
	ctx.JSON(http.StatusOK, response.SuccessResponse("login success", responseData))
	println("LoginHandler")
}

func (h *AuthAndler) OAuthHandler(ctx *gin.Context) {
	println("OAuthHandler")
}
