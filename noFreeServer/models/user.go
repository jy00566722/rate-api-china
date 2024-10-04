package models

import (
	"time"

	"gorm.io/gorm"
)

// MembershipLevel defines the type for membership levels
type MembershipLevel string

const (
	FreeMember     MembershipLevel = "free"
	HalfYearMember MembershipLevel = "half_year"
	YearMember     MembershipLevel = "year"
	LifetimeMember MembershipLevel = "lifetime"
)

type User struct {
	gorm.Model
	Email              string          `gorm:"type:varchar(255);uniqueIndex"`
	Password           string          `gorm:"type:varchar(255)"`
	MembershipLevel    MembershipLevel `gorm:"type:varchar(20);default:'free'"`
	MembershipExpireAt time.Time
	MaxDevices         int
	Devices            []Device `gorm:"foreignKey:UserID"`
}

type Device struct {
	gorm.Model
	UserID         uint   `gorm:"index"`
	Fingerprint    string `gorm:"type:varchar(255);uniqueIndex"`
	LastActiveTime time.Time
	User           User `gorm:"foreignKey:UserID"`
}

// IsValidMembership checks if the user's membership is still valid
func (u *User) IsValidMembership() bool {
	if u.MembershipLevel == FreeMember {
		return true
	}
	if u.MembershipLevel == LifetimeMember {
		return true
	}
	return time.Now().Before(u.MembershipExpireAt)
}

// CanAddDevice checks if the user can add more devices
func (u *User) CanAddDevice() bool {
	return len(u.Devices) < u.MaxDevices
}

// GetDeviceLimit returns the device limit based on membership level
func (u *User) GetDeviceLimit() int {
	switch u.MembershipLevel {
	case HalfYearMember:
		return 5
	case YearMember:
		return 10
	case LifetimeMember:
		return 999999 // practically unlimited
	default:
		return 1 // free users get 1 device
	}
}
