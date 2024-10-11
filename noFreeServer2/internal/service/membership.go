package service

import (
	"context"

	"mihu007/internal/model"
	"mihu007/internal/repository"
)

type MembershipService interface {
	GetMembershipInfo(ctx context.Context, userID uint) (*model.MemberLevelInfo, error)
	PurchaseMembership(ctx context.Context, userID uint, planID uint) (*model.Order, error)
	// GetMembershipPlans(ctx context.Context) ([]model.MembershipPlan, error)
}

type membershipService struct {
	membershipRepo repository.MembershipRepository
	deviceRepo     repository.DeviceRepository
	orderRepo      repository.OrderRepository
	// planRepo       repository.MembershipPlanRepository
}

func NewMembershipService(
	membershipRepo repository.MembershipRepository,
	deviceRepo repository.DeviceRepository,
	orderRepo repository.OrderRepository,
	// planRepo repository.MembershipPlanRepository,
) MembershipService {
	return &membershipService{
		membershipRepo: membershipRepo,
		deviceRepo:     deviceRepo,
		orderRepo:      orderRepo,
		// planRepo:       planRepo,
	}
}

func (s *membershipService) GetMembershipInfo(ctx context.Context, userID uint) (*model.MemberLevelInfo, error) {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	deviceCount, err := s.deviceRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	levelName := getLevelName(membership.Level)

	return &model.MemberLevelInfo{
		Level:          membership.Level,
		Name:           levelName,
		DurationMonths: membership.DurationMonths,
		DeviceLimit:    membership.DeviceLimit,
		DeviceCount:    deviceCount,
	}, nil
}

func (s *membershipService) PurchaseMembership(ctx context.Context, userID uint, planID uint) (*model.Order, error) {

	order := &model.Order{
		UserID: userID,
		Status: "pending",
	}

	return order, nil
}

func (s *membershipService) GetMembershipPlans(ctx context.Context) ([]model.MembershipPlan, error) {
	return s.planRepo.GetActivePlans(ctx)
}
