package service

import (
	"context"
	"mihu007/internal/model"
)

type OrderService interface {
	GetOrders(ctx context.Context, userID uint) ([]model.Order, error)
	GetOrder(ctx context.Context, userID uint, orderID string) (*model.Order, error)
	CreateOrder(ctx context.Context, userID uint, productID string, quantity int) (*model.Order, error)
	UpdateOrder(ctx context.Context, userID uint, orderID string, status string) (*model.Order, error)
	DeleteOrder(ctx context.Context, userID uint, orderID string) error
}
