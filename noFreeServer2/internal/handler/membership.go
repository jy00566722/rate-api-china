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

func (h *MembershipHandler) GetMembershipInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "user not authenticated"})
		return
	}

	info, err := h.membershipService.GetMembershipInfo(c, userID.(uint))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, info)
}

func (h *MembershipHandler) PurchaseMembership(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "user not authenticated"})
		return
	}

	var req model.MembershipPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	order, err := h.membershipService.PurchaseMembership(c, userID.(uint), req.PlanID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, order)
}

func (h *MembershipHandler) GetMembershipPlans(c *gin.Context) {
	plans, err := h.membershipService.GetMembershipPlans(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, plans)
}
