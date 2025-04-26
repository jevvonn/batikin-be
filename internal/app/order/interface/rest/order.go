package rest

import (
	"batikin-be/internal/app/order/usecase"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/infra/validator"
	"batikin-be/internal/middleware"
	"batikin-be/internal/models"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecaseItf
	validator    validator.ValidationService
}

func NewOrderHandler(
	router fiber.Router,
	orderUsecase usecase.OrderUsecaseItf,
	validator validator.ValidationService,
) {
	handler := &OrderHandler{orderUsecase, validator}
	router.Post("/orders", middleware.Authenticated, handler.CreateOrder)
	router.Get("/orders", middleware.Authenticated, handler.GetAllByUserId)
	router.Get("/orders/:id", middleware.Authenticated, handler.GetSpecific)
}

func (h *OrderHandler) GetAllByUserId(ctx *fiber.Ctx) error {
	orders, err := h.orderUsecase.GetAllByUserId(ctx)
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
			Data:    orders,
		},
	)
}

func (h *OrderHandler) GetSpecific(ctx *fiber.Ctx) error {
	order, err := h.orderUsecase.GetSpecific(ctx)
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
			Data:    order,
		},
	)
}

func (h *OrderHandler) CreateOrder(ctx *fiber.Ctx) error {
	var req dto.CreateOrderRequest
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

	orders, err := h.orderUsecase.Create(ctx, req)
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
			Data:    orders,
		},
	)
}
