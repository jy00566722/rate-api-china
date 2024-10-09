package service

import (
	"context"
	"time"

	"mihu007/internal/model"
	"mihu007/internal/repository"
)

type MembershipService interface {
	GetMembershipInfo(ctx context.Context, userID uint) (*model.MembershipInfo, error)
	PurchaseMembership(ctx context.Context, userID uint, planID uint) (*model.Order, error)
	GetMembershipPlans(ctx context.Context) ([]model.MembershipPlan, error)
}

type membershipService struct {
	membershipRepo repository.MembershipRepository
	deviceRepo     repository.DeviceRepository
	orderRepo      repository.OrderRepository
	planRepo       repository.MembershipPlanRepository
}

func NewMembershipService(
	membershipRepo repository.MembershipRepository,
	deviceRepo repository.DeviceRepository,
	orderRepo repository.OrderRepository,
	planRepo repository.MembershipPlanRepository,
) MembershipService {
	return &membershipService{
		membershipRepo: membershipRepo,
		deviceRepo:     deviceRepo,
		orderRepo:      orderRepo,
		planRepo:       planRepo,
	}
}

func (s *membershipService) GetMembershipInfo(ctx context.Context, userID uint) (*model.MembershipInfo, error) {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	deviceCount, err := s.deviceRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	levelName := getLevelName(membership.Level)

	return &model.MembershipInfo{
		Level:       membership.Level,
		LevelName:   levelName,
		ExpireAt:    membership.ExpireAt,
		DeviceLimit: membership.DeviceLimit,
		DeviceCount: deviceCount,
	}, nil
}

func (s *membershipService) PurchaseMembership(ctx context.Context, userID uint, planID uint) (*model.Order, error) {
	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	// 创建订单
	order := &model.Order{
		UserID:   userID,
		PlanID:   planID,
		Amount:   plan.Price,
		Status:   "pending",
		ExpireAt: time.Now().Add(5 * time.Minute), // 5分钟后过期
	}

	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *membershipService) GetMembershipPlans(ctx context.Context) ([]model.MembershipPlan, error) {
	return s.planRepo.GetActivePlans(ctx)
}
