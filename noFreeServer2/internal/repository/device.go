// internal/repository/device.go
package repository

import (
	"context"
	"errors"
	"mihu007/internal/model"
	"time"

	"gorm.io/gorm"
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

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Create(device).Error
}

func (r *deviceRepository) GetByID(ctx context.Context, id uint) (*model.Device, error) {
	var device model.Device
	if err := r.db.WithContext(ctx).First(&device, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) GetByUserID(ctx context.Context, userID uint) ([]model.Device, error) {
	var devices []model.Device
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepository) GetByFingerprint(ctx context.Context, fingerprint string) (*model.Device, error) {
	var device model.Device
	if err := r.db.WithContext(ctx).Where("fingerprint = ?", fingerprint).First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) UpdateLastActiveAt(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Device{}).Where("id = ?", id).Update("last_active_at", time.Now()).Error
}

func (r *deviceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Device{}, id).Error
}

func (r *deviceRepository) CountByUserID(ctx context.Context, userID uint) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Device{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
