package entity

import "github.com/google/uuid"

type ProductSizeVariant struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id,omitempty"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	Size  string  `gorm:"type:varchar(50);not null" json:"size,omitempty"`
	Price float64 `gorm:"type:decimal(10,2);not null" json:"price,omitempty"`
}
