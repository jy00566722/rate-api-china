package services

import (
	"errors"
	"noFree/models"
	"noFree/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrDeviceLimitReached = errors.New("device limit reached")
	ErrDeviceNotFound     = errors.New("device not found")
	ErrInvalidMembership  = errors.New("invalid membership")
)

type UserService interface {
	RegisterUser(email, password string) error
	AuthenticateUser(email, password string) (*models.User, error)
	UpdateMembership(userID uint, level models.MembershipLevel) error
	AddDevice(userID uint, fingerprint string) error
	RemoveDevice(userID uint, fingerprint string) error
	ValidateDevice(userID uint, fingerprint string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) RegisterUser(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:           email,
		Password:        string(hashedPassword),
		MembershipLevel: models.FreeMember,
		MaxDevices:      1,
	}

	return s.userRepo.CreateUser(user)
}

func (s *userService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

func (s *userService) UpdateMembership(userID uint, level models.MembershipLevel) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	user.MembershipLevel = level
	user.MaxDevices = user.GetDeviceLimit()

	switch level {
	case models.HalfYearMember:
		user.MembershipExpireAt = time.Now().AddDate(0, 6, 0)
	case models.YearMember:
		user.MembershipExpireAt = time.Now().AddDate(1, 0, 0)
	case models.LifetimeMember:
		user.MembershipExpireAt = time.Now().AddDate(100, 0, 0)
	}

	return s.userRepo.UpdateUser(user)
}

func (s *userService) AddDevice(userID uint, fingerprint string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if !user.CanAddDevice() {
		if err := s.userRepo.RemoveOldestDevice(userID); err != nil {
			return err
		}
	}

	return s.userRepo.AddDevice(userID, fingerprint)
}

func (s *userService) RemoveDevice(userID uint, fingerprint string) error {
	return s.userRepo.RemoveDevice(userID, fingerprint)
}

func (s *userService) ValidateDevice(userID uint, fingerprint string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if !user.IsValidMembership() {
		return ErrInvalidMembership
	}

	devices, err := s.userRepo.GetUserDevices(userID)
	if err != nil {
		return err
	}

	deviceFound := false
	for _, device := range devices {
		if device.Fingerprint == fingerprint {
			deviceFound = true
			break
		}
	}

	if !deviceFound {
		return ErrDeviceNotFound
	}

	return s.userRepo.UpdateDeviceLastActive(fingerprint)
}
