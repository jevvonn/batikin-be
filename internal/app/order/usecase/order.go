package usecase

import (
	"batikin-be/internal/app/order/repository"
	productRepo "batikin-be/internal/app/product/repository"
	"batikin-be/internal/constant"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/domain/entity"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderUsecaseItf interface {
	GetAllByUserId(ctx *fiber.Ctx) ([]entity.Order, error)
	GetSpecific(ctx *fiber.Ctx) (entity.Order, error)
	Create(ctx *fiber.Ctx, req dto.CreateOrderRequest) (dto.CreateOrderResponse, error)
}

type OrderUsecase struct {
	orderRepo   repository.OrderPostgreSQLItf
	productRepo productRepo.ProductPostgreSQLItf
}

func NewOrderUsecase(orderRepo repository.OrderPostgreSQLItf, productRepo productRepo.ProductPostgreSQLItf) OrderUsecaseItf {
	return &OrderUsecase{orderRepo, productRepo}
}

func (u *OrderUsecase) GetAllByUserId(ctx *fiber.Ctx) ([]entity.Order, error) {
	local := ctx.Locals("userId").(string)
	userId, err := uuid.Parse(local)
	if err != nil {
		return nil, err
	}

	return u.orderRepo.GetAllByUserId(userId)
}

func (u *OrderUsecase) GetSpecific(ctx *fiber.Ctx) (entity.Order, error) {
	param := ctx.Params("id")
	orderId, err := uuid.Parse(param)
	if err != nil {
		return entity.Order{}, err
	}

	local := ctx.Locals("userId").(string)
	userId, err := uuid.Parse(local)
	if err != nil {
		return entity.Order{}, err
	}

	return u.orderRepo.GetSpecific(entity.Order{
		ID:     orderId,
		UserID: userId,
	})
}

func (u *OrderUsecase) Create(ctx *fiber.Ctx, req dto.CreateOrderRequest) (dto.CreateOrderResponse, error) {
	local := ctx.Locals("userId").(string)
	userId, err := uuid.Parse(local)
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}

	productId, err := uuid.Parse(req.ProductID)
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}

	sizeVariantId, err := uuid.Parse(req.SizeVariantID)
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}

	product, err := u.productRepo.GetSpecific(entity.Product{
		ID: productId,
	})
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}

	size := entity.ProductSizeVariant{}
	for _, variant := range product.Sizes {
		if variant.ID == sizeVariantId {
			size = variant
			break
		}
	}

	if size.ID == uuid.Nil {
		return dto.CreateOrderResponse{}, fmt.Errorf("size variant not found")
	}

	totalPrice := size.Price * float64(req.Quantity)

	orderId := uuid.New()
	order := &entity.Order{
		ID:             orderId,
		UserID:         userId,
		ProductID:      productId,
		SizeVariantID:  sizeVariantId,
		Quantity:       req.Quantity,
		ProductionType: req.ProductionType,
		Address:        req.Address,
		TotalPrice:     totalPrice,
		Status:         constant.OrderPending,
	}

	err = u.orderRepo.Create(order)
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}

	orderDetail, err := u.orderRepo.GetSpecific(entity.Order{
		ID: orderId,
	})
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}

	return dto.CreateOrderResponse{
		OrderDetail: orderDetail,
	}, err
}
