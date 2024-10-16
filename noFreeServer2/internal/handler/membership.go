// File: internal/handler/membership.go
package handler

import (
	"mihu007/internal/model"
	"mihu007/internal/service"

	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	membershipService service.MembershipService
}

func NewMembershipHandler(membershipService service.MembershipService) *MembershipHandler {
	return &MembershipHandler{
		membershipService: membershipService,
	}
}

// Web端接口

// GetMembershipInfo 获取会员信息
// GET /api/v1/membership/info
func (h *MembershipHandler) GetMembershipInfo(c *gin.Context) {
	userID := getUserIDFromContext(c)
	info, err := h.membershipService.GetMembershipInfo(c, userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, info)
}

// GetMembershipPlans 获取会员计划列表
// GET /api/v1/membership/plans
func (h *MembershipHandler) GetMembershipPlans(c *gin.Context) {
	plans := model.MemberLevelInfos
	c.JSON(200, plans)
}

// PurchaseMembership 购买会员
// POST /api/v1/membership/purchase
func (h *MembershipHandler) PurchaseMembership(c *gin.Context) {
	userID := getUserIDFromContext(c)
	var req struct {
		Level uint `json:"level" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	order, err := h.membershipService.PurchaseMembership(c, userID, uint(req.Level))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, order)
}

// 插件端接口

// VerifyMembership 验证会员状态
// GET /api/v1/membership/verify
func (h *MembershipHandler) VerifyMembership(c *gin.Context) {
	userID := getUserIDFromContext(c)
	isValid, err := h.membershipService.VerifyMembership(c, userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"isValid": isValid})
}

// RegisterDevice 注册设备
// POST /api/v1/device/register
func (h *MembershipHandler) RegisterDevice(c *gin.Context) {
	userID := getUserIDFromContext(c)
	var req struct {
		Fingerprint string `json:"fingerprint" binding:"required"`
		ExtID       string `json:"extID" binding:"required"`
		UserAgent   string `json:"userAgent"`
		IPAddress   string `json:"ipAddress"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := h.membershipService.RegisterDevice(c, userID, req.Fingerprint, req.UserAgent, req.IPAddress)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Device registered successfully"})
}

func getUserIDFromContext(c *gin.Context) uint {
	// 从 JWT token 中获取用户 ID 的逻辑
	// 这里假设已经通过中间件设置了 userID
	userID, _ := c.Get("userID")
	return userID.(uint)
}
