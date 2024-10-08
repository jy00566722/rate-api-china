package model

import (
	"time"

	"gorm.io/gorm"
)

// 会员等级常量
const (
	MemberLevelFree      = 0
	MemberLevelNormal    = 1
	MemberLevelPremium   = 2
	MemberLevelUnlimited = 3
)

type Membership struct {
	gorm.Model
	UserID      uint      `gorm:"not null"`
	Level       int       `gorm:"not null"` // 0: free, 1: normal, 2: premium, 3: unlimited
	ExpireAt    time.Time `gorm:"not null"`
	DeviceLimit int       `gorm:"not null"`
}

type MembershipPlan struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Level       int     `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Duration    int     `gorm:"not null"` // 有效期（天）
	DeviceLimit int     `gorm:"not null"`
	Status      int     `gorm:"default:1"` // 0: disabled, 1: active
}

// MembershipPlanRequest 购买会员请求
type MembershipPlanRequest struct {
	PlanID uint `json:"plan_id" binding:"required"`
}

// MembershipInfo 会员信息响应
type MembershipInfo struct {
	Level       int       `json:"level"`
	LevelName   string    `json:"level_name"`
	ExpireAt    time.Time `json:"expire_at"`
	DeviceLimit int       `json:"device_limit"`
	DeviceCount int       `json:"device_count"`
}
