package migrations

import (
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/entity"
	"github.com/yaza-putu/online-test-bookandlink/internal/database"
	"gorm.io/gorm"
)

/// please replace or change &EntityName{}
/// AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes.
// It will change existing column’s type if its size, precision, nullable changed.
// It WON’T delete unused columns to protect your data.

func init() {
	database.MigrationRegister(func(db *gorm.DB) error {
		return db.AutoMigrate(&entity.FailedJob{})
	}, func(db *gorm.DB) error {
		return db.Migrator().DropTable(&entity.FailedJob{})
	})
}