package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID       uint
	MembershipID uint
	Membership   Membership
	Status       string
	Amount       float64
	ExpireAt     time.Time
}
