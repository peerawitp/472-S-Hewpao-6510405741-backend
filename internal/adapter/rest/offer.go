package rest

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/util"
)

type OfferHandler interface {
	CreateOffer(c *fiber.Ctx) error
	GetOfferDetailByOfferID(c *fiber.Ctx) error
}

type offerHandler struct {
	service usecase.OfferUsecase
}

func NewOfferHandler(service usecase.OfferUsecase) OfferHandler {
	return &offerHandler{
		service: service,
	}
}

func (o *offerHandler) CreateOffer(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	userID, _ := claims["id"].(string)

	req := dto.CreateOfferDTO{}
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

	err = o.service.CreateOffer(&req, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":   "Offer successfully created!",
		"offer-req": req,
	})
}

func (o *offerHandler) GetOfferDetailByOfferID(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	userID, _ := claims["id"].(string)

	param := c.Params("id")
	paramInt, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "offerID must be a number",
		})
	}

	offerDetail, err := o.service.GetOfferDetailByOfferID(paramInt, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(offerDetail)
}
