package handler

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"

	"github.com/mxxmstar/learning/verify_server/internal/service"
)

// 用户业务逻辑处理
// 用户相关的数据访问对象(DAO)初始化
// 定义用户相关的 HTTP 处理函数
// 用户数据验证等

type UserHandler struct {
	// userService 用户服务
	userService *service.UserService
	// emailExp 邮箱正则表达式
	emailExp *regexp.Regexp
	// passwordExp 密码正则表达式
	passwordExp *regexp.Regexp
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		emailExp:    regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, regexp.None),
		passwordExp: regexp.MustCompile(`^[a-zA-Z0-9_-]{6,20}$`, regexp.None),
	}
}

func (h *UserHandler) ProfileHandler(ctx *gin.Context) {

}

func (h *UserHandler) UpdateProfileHandler(ctx *gin.Context) {

}
