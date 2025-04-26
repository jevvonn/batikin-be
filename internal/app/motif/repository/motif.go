package repository

import (
	"batikin-be/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MotifPostgreSQLItf interface {
	GetAll() ([]entity.Motif, error)
	GetSpecific(motif entity.Motif) (entity.Motif, error)
	Create(motif *entity.Motif) error
}

type MotifPostgreSQL struct {
	db *gorm.DB
}

func NewMotifPostgreSQL(db *gorm.DB) MotifPostgreSQLItf {
	return &MotifPostgreSQL{db}
}

func (r *MotifPostgreSQL) GetAll() ([]entity.Motif, error) {
	var motifs []entity.Motif
	err := r.db.Preload("User").Find(&motifs).Error
	return motifs, err
}

func (r *MotifPostgreSQL) GetSpecific(motif entity.Motif) (entity.Motif, error) {
	var result entity.Motif
	err := r.db.Preload("User").First(&result, &motif).Error
	return result, err
}

func (r *MotifPostgreSQL) Create(motif *entity.Motif) error {
	motif.ID = uuid.New()
	return r.db.Create(&motif).Error
}
