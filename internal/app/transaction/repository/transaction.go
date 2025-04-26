package repository

import (
	"batikin-be/internal/domain/entity"

	"gorm.io/gorm"
)

type TransactionPostgreSQLItf interface {
	Create(transaction *entity.Transaction) error
	GetSpecific(transaction entity.Transaction) (entity.Transaction, error)
}

type TransactionPostgreSQL struct {
	db *gorm.DB
}

func NewTransactionPostgreSQL(db *gorm.DB) TransactionPostgreSQLItf {
	return &TransactionPostgreSQL{db}
}

func (r *TransactionPostgreSQL) GetSpecific(transaction entity.Transaction) (entity.Transaction, error) {
	var result entity.Transaction
	err := r.db.First(&result, &transaction).Error
	return result, err
}

func (r *TransactionPostgreSQL) Create(transaction *entity.Transaction) error {
	return r.db.Create(&transaction).Error
}
