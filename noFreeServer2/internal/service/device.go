package service

import (
	"context"
	"mihu007/internal/model"
)

// DeviceService 设备服务
type DeviceService interface {
	GetDevicesByUserID(ctx context.Context, userID string) ([]model.Device, error)
}
