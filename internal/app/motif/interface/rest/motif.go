package rest

import (
	"batikin-be/internal/app/motif/usecase"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/infra/validator"
	"batikin-be/internal/models"

	"github.com/gofiber/fiber/v2"
)

type MotifHandler struct {
	authUsecase usecase.MotifUsecaseItf
	validator   validator.ValidationService
}

func NewMotifHandler(
	router fiber.Router,
	authUsecase usecase.MotifUsecaseItf,
	validator validator.ValidationService,
) {
	handler := MotifHandler{authUsecase, validator}

	router.Post("/motif", handler.CreateMotif)
}

func (h *MotifHandler) CreateMotif(ctx *fiber.Ctx) error {
	var req dto.CreateMotifRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Invalid Request",
				Errors:  err.Error(),
			},
		)
	}

	err = h.validator.Validate(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			err.(*validator.ValidationError),
		)
	}

	res, err := h.authUsecase.Create(ctx, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Invalid Request",
				Errors:  err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		models.JSONResponseModel{
			Message: "Success",
			Data:    res,
		},
	)
}
