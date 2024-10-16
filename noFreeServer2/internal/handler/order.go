package handler

import (
	"mihu007/internal/model"
	"mihu007/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// Web端接口

// GetOrders 获取订单列表
// GET /api/v1/orders
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := getUserIDFromContext(c)
	orders, err := h.orderService.GetOrders(c, userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, orders)
}

// GetOrder 获取订单详情
// GET /api/v1/orders/:orderID
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID := getUserIDFromContext(c)
	orderID := c.Param("orderID")
	order, err := h.orderService.GetOrder(c, userID, orderID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, order)
}

// CreateOrder 创建订单
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		UserID    uint   `json:"user_id" binding:"required"`
		ProductID string `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(c, req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// 生成支付二维码
// POST /api/v1/orders/pay/:orderID
func (h *OrderHandler) GeneratePayCode(c *gin.Context) {
	orderID := c.Param("orderID")
	userID := getUserIDFromContext(c)
	order, err := h.orderService.GetOrder(c, userID, orderID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	} else if order.Status != model.OrderStatusPending {
		c.JSON(400, gin.H{"error": "Order status is not pending"})
		return
	}
}

// 支付回调
// POST /api/v1/orders/callback/:orderID
func (h *OrderHandler) Callback(c *gin.Context) {
	orderID := c.Param("orderID")
	userID := getUserIDFromContext(c)
	order, err := h.orderService.GetOrder(c, userID, orderID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	} else if order.Status != model.OrderStatusPending {
		c.JSON(400, gin.H{"error": "Order status is not pending"})
		return
	}
}
