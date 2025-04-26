package dto

type CreateMotifRequest struct {
	Name   string `json:"name"`
	Prompt string `json:"prompt" validate:"required"`
}

type CaptureMotifResponse struct {
	ImageURL string `json:"image_url"`
}
