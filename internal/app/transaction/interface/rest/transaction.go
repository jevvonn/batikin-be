package rest

import (
	"batikin-be/internal/app/transaction/usecase"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecaseItf
}

func NewTransactionHandler(
	router fiber.Router,
	transactionUsecase usecase.TransactionUsecaseItf,
) {
	handler := TransactionHandler{transactionUsecase}

	router.Post("/webhook/midtrans", handler.WebhookMidtrans)
}

func (h *TransactionHandler) WebhookMidtrans(ctx *fiber.Ctx) error {
	return h.transactionUsecase.WebhookMidtrans(ctx)
}
