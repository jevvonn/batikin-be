package rest

import (
	"batikin-be/internal/app/product/usecase"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/infra/validator"
	"batikin-be/internal/middleware"
	"batikin-be/internal/models"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecaseItf
	validator      validator.ValidationService
}

func NewProductHandler(
	router fiber.Router,
	productUsecase usecase.ProductUsecaseItf,
	validator validator.ValidationService,
) {
	handler := &ProductHandler{productUsecase, validator}

	router.Get("/products", handler.GetAll)
	router.Get("/products/:id", handler.GetSpecific)
	router.Post("/products/motif/:motifId", middleware.Authenticated, handler.Create)
}

func (h *ProductHandler) GetAll(ctx *fiber.Ctx) error {
	products, err := h.productUsecase.GetAll()
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
			Data:    products,
		},
	)
}

func (h *ProductHandler) GetSpecific(ctx *fiber.Ctx) error {
	product, err := h.productUsecase.GetSpecific(ctx)
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
			Data:    product,
		},
	)
}

func (h *ProductHandler) Create(ctx *fiber.Ctx) error {
	var req dto.CreateFromMotifProductRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
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

	product, err := h.productUsecase.CreateFromMotif(ctx, req)
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
			Data:    product,
		},
	)
}
