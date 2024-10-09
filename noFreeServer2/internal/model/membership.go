package model

import (
	"time"

	"gorm.io/gorm"
)

// MemberLevel 定义会员等级
type MemberLevel int

const (
	MemberLevelFree MemberLevel = iota
	MemberLevelNormal
	MemberLevelPremium
	MemberLevelUnlimited
)

// MemberLevelInfo 定义每个会员等级的详细信息
type MemberLevelInfo struct {
	Level          MemberLevel
	Name           string
	DeviceLimit    int
	DurationMonths int // -1 表示永久有效
}

// MemberLevelInfos 定义所有会员等级的详细信息
var MemberLevelInfos = map[MemberLevel]MemberLevelInfo{
	MemberLevelFree:      {Level: MemberLevelFree, Name: "免费会员", DeviceLimit: 1, DurationMonths: 1},
	MemberLevelNormal:    {Level: MemberLevelNormal, Name: "普通会员", DeviceLimit: 5, DurationMonths: 6},
	MemberLevelPremium:   {Level: MemberLevelPremium, Name: "高级会员", DeviceLimit: 10, DurationMonths: 12},
	MemberLevelUnlimited: {Level: MemberLevelUnlimited, Name: "永久会员", DeviceLimit: -1, DurationMonths: -1},
}

// Membership 定义用户的会员信息
type Membership struct {
	gorm.Model
	UserID      uint        `gorm:"not null"`
	Level       MemberLevel `gorm:"not null"`
	ExpireAt    time.Time   `gorm:"not null"`
	DeviceLimit int         `gorm:"not null"`
}

// MembershipPlan 定义会员计划
type MembershipPlan struct {
	gorm.Model
	Level  MemberLevel `gorm:"not null"`
	Price  float64     `gorm:"not null"`
	Status int         `gorm:"default:1"` // 0: disabled, 1: active
}

// MembershipPlanRequest 定义购买会员请求
type MembershipPlanRequest struct {
	PlanID uint `json:"plan_id" binding:"required"`
}

// MembershipInfo 定义会员信息响应
type MembershipInfo struct {
	Level       MemberLevel `json:"level"`
	LevelName   string      `json:"level_name"`
	ExpireAt    time.Time   `json:"expire_at"`
	DeviceLimit int         `json:"device_limit"`
	DeviceCount int         `json:"device_count"`
}
