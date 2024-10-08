// internal/repository/errors.go
package repository

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrDeviceNotFound         = errors.New("device not found")
	ErrMembershipNotFound     = errors.New("membership not found")
	ErrMembershipPlanNotFound = errors.New("membership plan not found")
)
