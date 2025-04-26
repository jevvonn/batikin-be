package usecase

import (
	"batikin-be/internal/app/order/repository"
	productRepo "batikin-be/internal/app/product/repository"
	transactionRepo "batikin-be/internal/app/transaction/repository"
	"batikin-be/internal/constant"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/domain/entity"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type OrderUsecaseItf interface {
	GetAllByUserId(ctx *fiber.Ctx) ([]entity.Order, error)
	GetSpecific(ctx *fiber.Ctx) (dto.GetOrderResponse, error)
	Create(ctx *fiber.Ctx, req dto.CreateOrderRequest) (dto.GetOrderResponse, error)
}

type OrderUsecase struct {
	orderRepo          repository.OrderPostgreSQLItf
	productRepo        productRepo.ProductPostgreSQLItf
	transactionRepo    transactionRepo.TransactionPostgreSQLItf
	midtransSnapClient snap.Client
}

func NewOrderUsecase(
	orderRepo repository.OrderPostgreSQLItf,
	productRepo productRepo.ProductPostgreSQLItf,
	midtransSnapClient snap.Client,
	transactionRepo transactionRepo.TransactionPostgreSQLItf,
) OrderUsecaseItf {
	return &OrderUsecase{orderRepo, productRepo, transactionRepo, midtransSnapClient}
}

func (u *OrderUsecase) GetAllByUserId(ctx *fiber.Ctx) ([]entity.Order, error) {
	local := ctx.Locals("userId").(string)
	userId, err := uuid.Parse(local)
	if err != nil {
		return nil, err
	}

	return u.orderRepo.GetAllByUserId(userId)
}

func (u *OrderUsecase) GetSpecific(ctx *fiber.Ctx) (dto.GetOrderResponse, error) {
	param := ctx.Params("id")
	orderId, err := uuid.Parse(param)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	local := ctx.Locals("userId").(string)
	userId, err := uuid.Parse(local)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	orderDetail, err := u.orderRepo.GetSpecific(entity.Order{
		ID:     orderId,
		UserID: userId,
	})
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	transactionDetail, err := u.transactionRepo.GetSpecific(entity.Transaction{
		ID: orderDetail.TransactionId,
	})
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	return dto.GetOrderResponse{
		OrderDetail:       orderDetail,
		TransactionDetail: transactionDetail,
	}, err
}

func (u *OrderUsecase) Create(ctx *fiber.Ctx, req dto.CreateOrderRequest) (dto.GetOrderResponse, error) {
	local := ctx.Locals("userId").(string)
	userId, err := uuid.Parse(local)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	productId, err := uuid.Parse(req.ProductID)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	sizeVariantId, err := uuid.Parse(req.SizeVariantID)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	product, err := u.productRepo.GetSpecific(entity.Product{
		ID: productId,
	})
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	size := entity.ProductSizeVariant{}
	for _, variant := range product.Sizes {
		if variant.ID == sizeVariantId {
			size = variant
			break
		}
	}

	if size.ID == uuid.Nil {
		return dto.GetOrderResponse{}, fmt.Errorf("size variant not found")
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
		return dto.GetOrderResponse{}, err
	}

	orderDetail, err := u.orderRepo.GetSpecific(entity.Order{
		ID: orderId,
	})
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	customerName := ctx.Locals("name").(string)
	customerEmail := ctx.Locals("email").(string)
	request := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			GrossAmt: int64(totalPrice),
			OrderID:  orderId.String(),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
			ShipAddr: &midtrans.CustomerAddress{
				Address: req.Address,
			},
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    orderDetail.ProductID.String(),
				Name:  orderDetail.Product.Name,
				Price: int64(orderDetail.SizeVariant.Price),
				Qty:   int32(orderDetail.Quantity),
			},
		},
	}

	url, err := u.midtransSnapClient.CreateTransactionUrl(request)

	transaction := &entity.Transaction{
		ID:         uuid.New(),
		OrderID:    orderId,
		Amount:     totalPrice,
		PaymentURL: url,
		Status:     constant.TransactionPending,
	}
	err = u.transactionRepo.Create(transaction)
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	transactionDetail, err := u.transactionRepo.GetSpecific(entity.Transaction{
		ID: transaction.ID,
	})
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	err = u.orderRepo.Update(&entity.Order{
		ID:            orderId,
		TransactionId: transaction.ID,
	})
	if err != nil {
		return dto.GetOrderResponse{}, err
	}

	return dto.GetOrderResponse{
		OrderDetail:       orderDetail,
		TransactionDetail: transactionDetail,
	}, err
}
