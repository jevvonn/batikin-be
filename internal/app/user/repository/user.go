package repository

import (
	"batikin-be/internal/domain/entity"

	"gorm.io/gorm"
)

type UserPostgreSQLItf interface {
	GetSpecificUser(user entity.User) (entity.User, error)
	CreateUser(user entity.User) error
}

type UserPostgreSQL struct {
	db *gorm.DB
}

func NewUserPostgreSQL(db *gorm.DB) UserPostgreSQLItf {
	return &UserPostgreSQL{db}
}

func (r *UserPostgreSQL) GetSpecificUser(user entity.User) (entity.User, error) {
	var result entity.User
	err := r.db.First(&result, &user).Error
	return result, err
}

func (r *UserPostgreSQL) CreateUser(user entity.User) error {
	return r.db.Create(&user).Error
}
