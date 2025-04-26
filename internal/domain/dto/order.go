package dto

import "batikin-be/internal/domain/entity"

type CreateOrderRequest struct {
	ProductID     string `json:"product_id,omitempty" validate:"required"`
	SizeVariantID string `json:"product_size_variant_id,omitempty" validate:"required"`

	Quantity       int    `json:"quantity,omitempty" validate:"required"`
	ProductionType string `json:"production_type,omitempty" validate:"required,oneof=batik_tulis batik_cetak"`
	Address        string `json:"address,omitempty" validate:"required"`
}

type GetOrderResponse struct {
	OrderDetail       entity.Order       `json:"order_detail,omitempty"`
	TransactionDetail entity.Transaction `json:"payment_details,omitempty"`
}
