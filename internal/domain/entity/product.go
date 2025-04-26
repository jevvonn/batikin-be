package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`
	Name      string    `gorm:"type:text;not null" json:"name,omitempty"`
	ImageURL  string    `gorm:"type:varchar(255);not null" json:"image_url,omitempty"`
	ClothType string    `gorm:"type:varchar(50);not null" json:"cloth_type,omitempty"`
	MotifID   uuid.UUID `gorm:"type:uuid;not null" json:"motif_id,omitempty"`

	Sizes []ProductSizeVariant `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sizes,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
