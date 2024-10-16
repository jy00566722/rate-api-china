package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	gorm.Model
	UserID        uint
	MemberLevel   MemberLevel
	Amount        float64
	Status        OrderStatus
	PaymentMethod string
	//支付订单号
	PaymentOrderID string    `gorm:"type:varchar(100);not null"`
	PaymentTime    time.Time `gorm:"not null"`
}
