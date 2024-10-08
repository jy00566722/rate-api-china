// internal/repository/membership.go
package repository

import (
	"context"
	"errors"
	"mihu007/internal/model"

	"gorm.io/gorm"
)

type MembershipRepository interface {
	Create(ctx context.Context, membership *model.Membership) error
	GetByUserID(ctx context.Context, userID uint) (*model.Membership, error)
	Update(ctx context.Context, membership *model.Membership) error
	GetActiveMembershipPlan(ctx context.Context, level int) (*model.MembershipPlan, error)
	GetAllActiveMembershipPlans(ctx context.Context) ([]model.MembershipPlan, error)
}

type membershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipRepository{db: db}
}

func (r *membershipRepository) Create(ctx context.Context, membership *model.Membership) error {
	return r.db.WithContext(ctx).Create(membership).Error
}

func (r *membershipRepository) GetByUserID(ctx context.Context, userID uint) (*model.Membership, error) {
	var membership model.Membership
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&membership).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMembershipNotFound
		}
		return nil, err
	}
	return &membership, nil
}

func (r *membershipRepository) Update(ctx context.Context, membership *model.Membership) error {
	return r.db.WithContext(ctx).Save(membership).Error
}

func (r *membershipRepository) GetActiveMembershipPlan(ctx context.Context, level int) (*model.MembershipPlan, error) {
	var plan model.MembershipPlan
	if err := r.db.WithContext(ctx).Where("level = ? AND status = 1", level).First(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMembershipPlanNotFound
		}
		return nil, err
	}
	return &plan, nil
}

func (r *membershipRepository) GetAllActiveMembershipPlans(ctx context.Context) ([]model.MembershipPlan, error) {
	var plans []model.MembershipPlan
	if err := r.db.WithContext(ctx).Where("status = 1").Order("level").Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}
