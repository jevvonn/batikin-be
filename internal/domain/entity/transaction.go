package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID      uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`
	OrderID uuid.UUID `gorm:"type:uuid;not null" json:"order_id,omitempty"`
	Order   Order     `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	Amount     float64 `gorm:"type:decimal(10,2);not null" json:"amount,omitempty"`
	PaymentURL string  `gorm:"type:varchar(255);not null" json:"payment_url,omitempty"`
	Status     string  `gorm:"type:varchar(50);not null" json:"status,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
