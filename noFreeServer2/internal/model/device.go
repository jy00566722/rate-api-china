package model

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	UserID       uint      `gorm:"not null"`
	Fingerprint  string    `gorm:"unique;not null"`
	LastActiveAt time.Time `gorm:"not null"`
	UserAgent    string
	IPAddress    string
	Status       int `gorm:"default:1"` // 0: disabled, 1: active
}

type DeviceRegisterRequest struct {
	Fingerprint string `json:"fingerprint" binding:"required"`
}
