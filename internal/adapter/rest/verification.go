package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"

	"github.com/hewpao/hewpao-backend/util"
)

type VerificationHandler interface {
	VerifyWithKYC(c *fiber.Ctx) error
	GetVerificationInfo(c *fiber.Ctx) error
	UpdateVerificationInfo(c *fiber.Ctx) error
}

type verificationHandler struct {
	service usecase.VerificationUsecase
}

func NewVerificationHandler(service usecase.VerificationUsecase) VerificationHandler {
	return &verificationHandler{
		service: service,
	}
}

func (v *verificationHandler) VerifyWithKYC(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("A valid authorization is required")
	}

	fileHeaders, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	fileReaders, files, err := util.FileManage(fileHeaders, "cardImages", 1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	err = v.service.VerifyWithKYC(fileReaders[0], files[0], token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "card uploaded",
	})
}

func (v *verificationHandler) GetVerificationInfo(c *fiber.Ctx) error {
	userEmail := c.Params("email")
	if userEmail == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Email param is missing")
	}

	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("A valid authorization is required")
	}

	information := dto.GetUserVerificationDTO{}

	err := v.service.GetVerificationInfo(userEmail, token, &information)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message":    "Get verifcation Info success",
		"informaion": information,
	})
}

func (v *verificationHandler) UpdateVerificationInfo(c *fiber.Ctx) error {
	userEmail := c.Params("email")
	if userEmail == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Email param is missing")
	}

	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("A valid authorization is required")
	}

	req := dto.UpdateUserVerificationDTO{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = v.service.UpdateVerificationInfo(&req, userEmail, token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "IsVerified has been set",
	})
}
