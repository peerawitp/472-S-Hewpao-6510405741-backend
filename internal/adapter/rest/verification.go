package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
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
	claims := c.Locals("user").(jwt.MapClaims)
	userEmail, _ := claims["email"].(string)

	req := dto.VerifyWithKYCDTO{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "some request fields are required",
		})
	}
	validationErr := util.ValidateStruct(req)
	if validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationErr.Error,
		})
	}

	fileHeaders, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fileReaders, files, err := util.FileManage(fileHeaders, "card-image", 1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = v.service.VerifyWithKYC(fileReaders[0], files[0], userEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "card uploaded",
	})
}

func (v *verificationHandler) GetVerificationInfo(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	instructorEmail, _ := claims["email"].(string)

	userEmail := c.Params("email")
	if userEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email param is missing",
		})
	}

	information := dto.GetUserVerificationDTO{}

	err := v.service.GetVerificationInfo(userEmail, instructorEmail, &information)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":    "Get verifcation Info success",
		"informaion": information,
	})
}

func (v *verificationHandler) UpdateVerificationInfo(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	instructorEmail, _ := claims["email"].(string)

	userEmail := c.Params("email")
	if userEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email param is missing",
		})
	}

	req := dto.UpdateUserVerificationDTO{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = v.service.UpdateVerificationInfo(&req, userEmail, instructorEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "IsVerified has been set",
	})
}
