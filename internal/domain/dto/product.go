package dto

type CreateFromMotifProductRequest struct {
	ClothType string `json:"cloth_type,omitempty" validate:"required,oneof=kemeja outer kain"`
}
