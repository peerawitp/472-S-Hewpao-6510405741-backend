package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/util"
)

type CheckoutHandler interface {
	CheckoutWithPaymentGateway(c *fiber.Ctx) error
}

type checkoutHandler struct {
	service usecase.CheckoutUsecase
}

func NewCheckoutHandler(service usecase.CheckoutUsecase) CheckoutHandler {
	return &checkoutHandler{
		service: service,
	}
}

func (ch *checkoutHandler) CheckoutWithPaymentGateway(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	userID, _ := claims["id"].(string)

	var req dto.CheckoutRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	payment, err := ch.service.CheckoutWithPaymentGateway(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(payment)
}
