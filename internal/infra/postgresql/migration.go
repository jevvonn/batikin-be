package postgresql

import (
	"batikin-be/internal/domain/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, command string) {
	migrator := db.Migrator()
	tables := []any{
		&entity.User{},
		&entity.Motif{},
		&entity.Product{},
		&entity.ProductSizeVariant{},
		&entity.Order{},
		&entity.Transaction{},
	}

	var err error
	if command == "up" {
		err = migrator.AutoMigrate(tables...)
	}

	if command == "down" {
		err = migrator.DropTable(tables...)
	}

	if err != nil {
		panic(err)
	}
}
