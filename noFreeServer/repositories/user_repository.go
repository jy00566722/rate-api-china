package repositories

import (
	"noFree/models"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
	AddDevice(userID uint, fingerprint string) error
	RemoveDevice(userID uint, fingerprint string) error
	GetUserDevices(userID uint) ([]models.Device, error)
	UpdateDeviceLastActive(fingerprint string) error
	RemoveOldestDevice(userID uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Devices").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Devices").First(&user, id).Error
	return &user, err
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) AddDevice(userID uint, fingerprint string) error {
	device := models.Device{
		UserID:         userID,
		Fingerprint:    fingerprint,
		LastActiveTime: time.Now(),
	}
	return r.db.Create(&device).Error
}

func (r *userRepository) RemoveDevice(userID uint, fingerprint string) error {
	return r.db.Where("user_id = ? AND fingerprint = ?", userID, fingerprint).Delete(&models.Device{}).Error
}

func (r *userRepository) GetUserDevices(userID uint) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.Where("user_id = ?", userID).Find(&devices).Error
	return devices, err
}

func (r *userRepository) UpdateDeviceLastActive(fingerprint string) error {
	return r.db.Model(&models.Device{}).
		Where("fingerprint = ?", fingerprint).
		Update("last_active_time", time.Now()).Error
}

func (r *userRepository) RemoveOldestDevice(userID uint) error {
	var oldestDevice models.Device
	err := r.db.Where("user_id = ?", userID).
		Order("last_active_time asc").
		First(&oldestDevice).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&oldestDevice).Error
}
