package service

import (
	"context"

	"mihu007/internal/model"
)

type MembershipService interface {
	GetMembershipInfo(ctx context.Context, userID uint) (*model.MemberLevelInfo, error)
	PurchaseMembership(ctx context.Context, userID uint, planID uint) (*model.Order, error)
	VerifyMembership(ctx context.Context, userID uint) (bool, error)
	RegisterDevice(ctx context.Context, userID uint, fingerprint, userAgent, ipAddress string) error
}
