package dto

type CreateMotifRequest struct {
	Name   string `json:"name"`
	Prompt string `json:"prompt" validate:"required"`
}
