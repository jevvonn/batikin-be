package usecase

import (
	"batikin-be/config"
	orderRepo "batikin-be/internal/app/order/repository"
	"batikin-be/internal/app/transaction/repository"
	"batikin-be/internal/constant"
	"batikin-be/internal/domain/entity"
	"batikin-be/internal/helper"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionUsecaseItf interface {
	WebhookMidtrans(ctx *fiber.Ctx) error
}

type TransactionUsecase struct {
	transactionRepo repository.TransactionPostgreSQLItf
	orderRepo       orderRepo.OrderPostgreSQLItf
}

func NewTransactionUsecase(
	transactionRepo repository.TransactionPostgreSQLItf,
	orderRepo orderRepo.OrderPostgreSQLItf,
) TransactionUsecaseItf {
	return &TransactionUsecase{transactionRepo, orderRepo}
}

type MidtransNotification struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
}

func (u *TransactionUsecase) WebhookMidtrans(ctx *fiber.Ctx) error {
	var notification MidtransNotification

	if err := ctx.BodyParser(&notification); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	serverKey := config.Load().MidtransServerKey
	expectedSignature := helper.GenerateSignature(notification.OrderID, notification.StatusCode, notification.GrossAmount, serverKey)

	if notification.SignatureKey != expectedSignature {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid signature",
		})
	}

	orderId, err := uuid.Parse(notification.OrderID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid order ID",
		})
	}

	order, err := u.orderRepo.GetSpecific(entity.Order{
		ID: orderId,
	})
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Invalid order ID",
		})
	}

	updateTransaction := &entity.Transaction{
		ID: order.TransactionId,
	}

	updateOrder := &entity.Order{
		ID: order.ID,
	}

	// Process transaction status
	switch notification.TransactionStatus {
	case "capture", "settlement":
		log.Println("Payment success for OrderID:", notification.OrderID)
		updateTransaction.Status = constant.TransactionDone
		updateOrder.Status = constant.OrderProcess
	case "deny", "cancel", "expire":
		updateTransaction.Status = constant.TransactionCanceled
		updateOrder.Status = constant.OrderCanceled
		log.Println("Payment failed for OrderID:", notification.OrderID)
	case "pending":
		updateTransaction.Status = constant.TransactionPending
	default:
		updateTransaction.Status = constant.TransactionCanceled
		updateOrder.Status = constant.OrderCanceled
		log.Println("Unhandled transaction status:", notification.TransactionStatus)
	}

	err = u.orderRepo.Update(updateOrder)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update order",
		})
	}

	err = u.transactionRepo.Update(updateTransaction)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{

			"message": "Failed to update transaction",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Transaction updated successfully",
	})
}
