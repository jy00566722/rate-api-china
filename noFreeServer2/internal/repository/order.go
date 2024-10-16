package repository

import (
	"context"
	"mihu007/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByUserID(ctx context.Context, userID uint) ([]model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}
