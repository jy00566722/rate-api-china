package model

import (
	"time"

	"gorm.io/gorm"
)

type MemberLevel int

const (
	MemberLevelFree MemberLevel = iota
	MemberLevelNormal
	MemberLevelPremium
	MemberLevelUnlimited
)

var MemberLevelInfos = map[MemberLevel]MemberLevelInfo{
	MemberLevelFree:      {Level: MemberLevelFree, Name: "免费会员", DeviceLimit: 1, DurationMonths: 1, Price: 0},
	MemberLevelNormal:    {Level: MemberLevelNormal, Name: "普通会员", DeviceLimit: 5, DurationMonths: 6, Price: 19},
	MemberLevelPremium:   {Level: MemberLevelPremium, Name: "高级会员", DeviceLimit: 10, DurationMonths: 12, Price: 29},
	MemberLevelUnlimited: {Level: MemberLevelUnlimited, Name: "永久会员", DeviceLimit: -1, DurationMonths: -1, Price: 59},
}

type MemberLevelInfo struct {
	Level          MemberLevel
	Name           string
	DeviceLimit    int
	DurationMonths int
	Price          float64
}

type Membership struct {
	gorm.Model
	UserID   uint
	Level    MemberLevel
	ExpireAt time.Time
	IsActive bool
}
