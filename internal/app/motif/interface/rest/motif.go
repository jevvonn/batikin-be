package rest

import (
	"batikin-be/internal/app/motif/usecase"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/infra/validator"
	"batikin-be/internal/middleware"
	"batikin-be/internal/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	router.Get("/motif", handler.GetAllMotif)
	router.Get("/motif/:id", handler.GetSpecific)
	router.Post("/motif", middleware.Authenticated, handler.CreateMotif)
}

func (h *MotifHandler) GetAllMotif(ctx *fiber.Ctx) error {
	motifs, err := h.authUsecase.GetAll()
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
			Data:    motifs,
		},
	)
}

func (h *MotifHandler) GetSpecific(ctx *fiber.Ctx) error {
	param := ctx.Params("id")
	_, err := uuid.Parse(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			models.JSONResponseModel{
				Message: "Invalid Request",
				Errors:  fmt.Errorf("not a valid id").Error(),
			},
		)
	}

	motif, err := h.authUsecase.GetSpecific(ctx)
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
			Data:    motif,
		},
	)
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
