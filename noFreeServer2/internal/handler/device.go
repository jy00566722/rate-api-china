package handler

import (
	"mihu007/internal/service"
	"mihu007/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	deviceService service.DeviceService
}

func NewDeviceHandler(deviceService service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// GetUserDevices 获取用户设备列表
func (h *DeviceHandler) GetUserDevices(c *gin.Context) {
	userID := c.Param("user_id")
	devices, err := h.deviceService.GetDevicesByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user devices"})
		return
	}

	c.JSON(http.StatusOK, devices)
}

type PluginLoginHandler struct {
	userService service.UserService
	jwtUtil     utils.JWTUtil
}

func NewPluginLoginHandler(userService service.UserService, jwtUtil utils.JWTUtil) *PluginLoginHandler {
	return &PluginLoginHandler{
		userService: userService,
		jwtUtil:     jwtUtil,
	}
}

// GenerateTemporaryToken 生成临时token
func (h *PluginLoginHandler) GenerateTemporaryToken(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	tempToken, err := h.jwtUtil.GenerateTemporaryToken(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate temporary token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"temp_token": tempToken})
}

// ExchangeTokenForJWT 用临时token换取JWT
func (h *PluginLoginHandler) ExchangeTokenForJWT(c *gin.Context) {
	tempToken := c.Query("temp_token")
	if tempToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Temporary token is required"})
		return
	}

	claims, err := h.jwtUtil.ValidateTemporaryToken(tempToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid temporary token"})
		return
	}

	jwtToken, err := h.jwtUtil.GenerateToken(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt_token": jwtToken})
}
