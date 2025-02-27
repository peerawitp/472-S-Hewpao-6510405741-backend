package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
)

type TransactionHandler interface {
	CreateTransaction(c *fiber.Ctx) error
	GetTransactionByID(c *fiber.Ctx) error
}

type transactionHandler struct {
	service usecase.TransactionService
}

func NewTransactionHandler(service usecase.TransactionService) TransactionHandler {
	return &transactionHandler{
		service: service,
	}
}

func (th *transactionHandler) CreateTransaction(c *fiber.Ctx) error {
	req := dto.CreateTransactionRequestDTO{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	transaction, err := th.service.CreateTransaction(req.UserID, req.Amount, req.Currency, req.Type)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(transaction)
}

func (th *transactionHandler) GetTransactionByID(c *fiber.Ctx) error {
	transactionID := c.Params("id")

	transaction, err := th.service.GetTransactionByID(c.Context(), transactionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	return c.JSON(transaction)
}
