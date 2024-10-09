package handler

import (
	"mihu007/internal/model"
	"mihu007/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(c, req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.userService.Login(c, req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *UserHandler) SendPasswordResetVerifyCode(c *gin.Context) {
	var req model.PasswordResetVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.SendPasswordResetVerifyCode(c, req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "verify code sent"})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req model.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.ResetPassword(c, req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "password reset successfully"})
}

// 需要新增或修改的处理函数
func (h *UserHandler) WechatLogin(c *gin.Context)
func (h *UserHandler) WechatRegister(c *gin.Context)
func (h *UserHandler) GenerateCaptcha(c *gin.Context)
func (h *UserHandler) VerifyCaptcha(c *gin.Context)
