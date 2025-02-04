package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/util"
)

type AuthHandler interface {
	LoginWithCredentials(c *fiber.Ctx) error
	LoginWithOAuth(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
}

type authHandler struct {
	service usecase.AuthUsecase
}

func NewAuthHandler(service usecase.AuthUsecase) AuthHandler {
	return &authHandler{
		service: service,
	}
}

func (a *authHandler) LoginWithCredentials(c *fiber.Ctx) error {
	var req dto.LoginWithCredentialsRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	user, err := a.service.LoginWithCredentials(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user)
}

func (a *authHandler) LoginWithOAuth(c *fiber.Ctx) error {
	var req dto.LoginWithOAuthRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	user, err := a.service.LoginWithOAuth(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(user)
}

func (a *authHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterUserRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	err := a.service.Register(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}
