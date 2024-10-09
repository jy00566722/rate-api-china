package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique"`
	Phone        string `gorm:"unique"`
	WechatID     string `gorm:"unique"`
	Password     string `gorm:"not null"`
	RegisterType string `gorm:"not null"`  // phone, email, wechat
	Status       int    `gorm:"default:1"` // 0: disabled, 1: active
	Nickname     string
	AvatarURL    string
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 可以是用户名、邮箱或手机号
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	VerifyCode string `json:"verify_code" binding:"required"`
}

// PasswordResetRequest 密码重置请求
type PasswordResetRequest struct {
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	VerifyCode  string `json:"verify_code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// PasswordResetVerifyRequest 获取密码重置验证码请求
type PasswordResetVerifyRequest struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type WechatInfo struct {
	OpenID    string
	UnionID   string
	Nickname  string
	AvatarURL string
}
