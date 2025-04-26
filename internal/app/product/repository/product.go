package repository

import (
	"batikin-be/internal/domain/entity"

	"gorm.io/gorm"
)

type ProductPostgreSQLItf interface {
	GetAll() ([]entity.Product, error)
	GetSpecific(product entity.Product) (entity.Product, error)
	Create(product *entity.Product) error
}

type ProductPostgreSQL struct {
	db *gorm.DB
}

func NewProductPostgreSQL(db *gorm.DB) ProductPostgreSQLItf {
	return &ProductPostgreSQL{db}
}

func (r *ProductPostgreSQL) GetAll() ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.Preload("Sizes").Find(&products).Error
	return products, err
}

func (r *ProductPostgreSQL) GetSpecific(product entity.Product) (entity.Product, error) {
	var result entity.Product
	err := r.db.Preload("Sizes").First(&result, &product).Error
	return result, err
}

func (r *ProductPostgreSQL) Create(product *entity.Product) error {
	return r.db.Create(&product).Error
}
