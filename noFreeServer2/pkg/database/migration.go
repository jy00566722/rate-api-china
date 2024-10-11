// pkg/database/migration.go
package database

import (
	"mihu007/internal/model"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Device{},
		&model.Membership{},
		&model.Order{},
	)
}
