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
	Price          float64
}

// MemberLevelInfos 定义所有会员等级的详细信息
// 1. 免费会员：1 个设备，1 个月
// 2. 普通会员：5 个设备，6 个月
// 3. 高级会员：10 个设备，12 个月
// 4. 永久会员：不限设备，不限时间
var MemberLevelInfos = map[MemberLevel]MemberLevelInfo{
	MemberLevelFree:      {Level: MemberLevelFree, Name: "免费会员", DeviceLimit: 1, DurationMonths: 1, Price: 0},
	MemberLevelNormal:    {Level: MemberLevelNormal, Name: "普通会员", DeviceLimit: 5, DurationMonths: 6, Price: 19},
	MemberLevelPremium:   {Level: MemberLevelPremium, Name: "高级会员", DeviceLimit: 10, DurationMonths: 12, Price: 29},
	MemberLevelUnlimited: {Level: MemberLevelUnlimited, Name: "永久会员", DeviceLimit: -1, DurationMonths: -1, Price: 59},
}

// Membership 定义用户的会员信息
type Membership struct {
	gorm.Model
	UserID      uint        `gorm:"not null"`
	Level       MemberLevel `gorm:"not null"`
	ExpireAt    time.Time   `gorm:"not null"`
	DeviceLimit int         `gorm:"not null"`
	IsActive    bool        `gorm:"default:true"`
}

// MembershipPlanRequest 定义购买会员请求
type MembershipPlanRequest struct {
	PlanID uint `json:"plan_id" binding:"required"`
}
