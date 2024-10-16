// internal/repository/device.go
package repository

import (
	"context"
	"mihu007/internal/model"
)

type DeviceRepository interface {
	Create(ctx context.Context, device *model.Device) error
	GetByID(ctx context.Context, id uint) (*model.Device, error)
	GetByUserID(ctx context.Context, userID uint) ([]model.Device, error)
	GetByFingerprint(ctx context.Context, fingerprint string) (*model.Device, error)
	UpdateLastActiveAt(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
	CountByUserID(ctx context.Context, userID uint) (int, error)
}
