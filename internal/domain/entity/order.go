package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID     uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id,omitempty"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id,omitempty"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product,omitempty"`

	SizeVariantID uuid.UUID          `gorm:"type:uuid;not null" json:"product_size_variant_id,omitempty"`
	SizeVariant   ProductSizeVariant `gorm:"foreignKey:SizeVariantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product_size_variant,omitempty"`

	Quantity       int     `gorm:"type:int;not null" json:"quantity,omitempty"`
	ProductionType string  `gorm:"type:varchar(50);not null" json:"production_type,omitempty"`
	TotalPrice     float64 `gorm:"type:decimal(10,2);not null" json:"total_price,omitempty"`
	Address        string  `gorm:"type:text;not null" json:"address,omitempty"`
	Status         string  `gorm:"type:varchar(50);not null" json:"status,omitempty"`

	TransactionId uuid.UUID `gorm:"type:uuid" json:"transaction_id,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
