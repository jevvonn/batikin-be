package entity

import (
	"time"

	"github.com/google/uuid"
)

type Motif struct {
	ID       uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`
	Name     string    `gorm:"type:text;not null" json:"name,omitempty"`
	Prompt   string    `gorm:"type:text;not null" json:"prompt,omitempty"`
	ImageURL string    `gorm:"type:varchar(255);not null" json:"image_url,omitempty"`

	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id,omitempty"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
