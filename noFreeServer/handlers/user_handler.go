package handlers

import (
	"net/http"
	"noFree/config"
	"noFree/services"
	"noFree/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
	config      *config.Config
}

func NewUserHandler(userService services.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		config:      cfg,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.RegisterUser(req.Email, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required"`
		Fingerprint string `json:"fingerprint" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := h.userService.AddDevice(user.ID, req.Fingerprint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register device"})
		return
	}

	token, err := utils.GenerateJWT(h.config, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"email":           user.Email,
			"membershipLevel": user.MembershipLevel,
			"expireAt":        user.MembershipExpireAt,
		},
	})
}
func (h *UserHandler) ValidateDevice(c *gin.Context) {
	userID := c.GetUint("userID")
	fingerprint := c.GetString("fingerprint")

	if err := h.userService.ValidateDevice(userID, fingerprint); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "valid"})
}

func (h *UserHandler) RemoveDevice(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		Fingerprint string `json:"fingerprint" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.RemoveDevice(userID, req.Fingerprint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device removed successfully"})
}
