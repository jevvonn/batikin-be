package repository

import (
	"batikin-be/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderPostgreSQLItf interface {
	GetAllByUserId(userId uuid.UUID) ([]entity.Order, error)
	GetSpecific(order entity.Order) (entity.Order, error)
	Create(order *entity.Order) error
}

type OrderPostgreSQL struct {
	db *gorm.DB
}

func NewOrderPostgreSQL(db *gorm.DB) OrderPostgreSQLItf {
	return &OrderPostgreSQL{db}
}

func (r *OrderPostgreSQL) GetAllByUserId(userId uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.Preload("Product").Preload("User").Preload("SizeVariant").Where("user_id = ?", userId).Find(&orders).Error
	return orders, err
}

func (r *OrderPostgreSQL) GetSpecific(order entity.Order) (entity.Order, error) {
	var result entity.Order
	err := r.db.Preload("Product").Preload("User").Preload("SizeVariant").First(&result, &order).Error
	return result, err
}

func (r *OrderPostgreSQL) Create(order *entity.Order) error {
	return r.db.Create(&order).Error
}
