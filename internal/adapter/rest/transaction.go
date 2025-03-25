package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
)

type TransactionHandler interface {
	CreateTransaction(c *fiber.Ctx) error
	GetTransactionByID(c *fiber.Ctx) error
	GetTransactionByUserID(c *fiber.Ctx) error
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

	transaction, err := th.service.CreateTransaction(req.UserID, req.Amount, req.Currency)
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

	res := &dto.GetTransactionResponseDTO{
		ID:                  transaction.ID,
		UserID:              transaction.UserID,
		Amount:              transaction.Amount,
		Currency:            transaction.Currency,
		ThirdPartyGateway:   transaction.ThirdPartyGateway,
		ThirdPartyPaymentID: transaction.ThirdPartyPaymentID,
		ProductRequestID:    transaction.ProductRequestID,
		Status:              transaction.Status,
		CreatedAt:           transaction.CreatedAt,
		UpdatedAt:           transaction.UpdatedAt,
	}

	return c.JSON(res)
}

func (th *transactionHandler) GetTransactionByUserID(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	userID, _ := claims["id"].(string)

	transactions, err := th.service.GetTransactionsByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	res := make([]*dto.GetTransactionResponseDTO, 0)

	for _, transaction := range transactions {
		res = append(res, &dto.GetTransactionResponseDTO{
			ID:                  transaction.ID,
			UserID:              transaction.UserID,
			Amount:              transaction.Amount,
			Currency:            transaction.Currency,
			ThirdPartyGateway:   transaction.ThirdPartyGateway,
			ThirdPartyPaymentID: transaction.ThirdPartyPaymentID,
			ProductRequestID:    transaction.ProductRequestID,
			Status:              transaction.Status,
			CreatedAt:           transaction.CreatedAt,
			UpdatedAt:           transaction.UpdatedAt,
		})
	}

	return c.JSON(res)
}
