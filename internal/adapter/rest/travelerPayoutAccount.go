package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/util"
)

type TravelerPayoutAccountHandler interface {
	CreateTravelerPayoutAccount(c *fiber.Ctx) error
	GetAccountsByUserID(c *fiber.Ctx) error
	GetAllAvailableBank(c *fiber.Ctx) error
}

type travelerPayoutAccountHandler struct {
	service usecase.TravelerPayoutAccountUsecase
}

func NewTravelerPayoutAccountHandler(service usecase.TravelerPayoutAccountUsecase) TravelerPayoutAccountHandler {
	return &travelerPayoutAccountHandler{
		service: service,
	}
}

func (h *travelerPayoutAccountHandler) CreateTravelerPayoutAccount(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	userID, _ := claims["id"].(string)

	req := dto.CreateTravelerPayoutAccountRequestDTO{}

	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	validationErr := util.ValidateStruct(req)
	if validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationErr.Error,
		})
	}

	err = h.service.CreateTravelerPayoutAccount(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Traveler payout account created successfully",
	})
}

func (h *travelerPayoutAccountHandler) GetAccountsByUserID(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	accounts, err := h.service.GetAccountsByUserID(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(accounts)
}

func (h *travelerPayoutAccountHandler) GetAllAvailableBank(c *fiber.Ctx) error {
	banks, err := h.service.GetAllAvailableBank()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(banks)
}
