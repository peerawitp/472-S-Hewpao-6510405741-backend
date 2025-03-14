package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/util"
)

type UserHandler interface {
	GetMyProfile(c *fiber.Ctx) error
	EditMyProfile(c *fiber.Ctx) error
}

type userHandler struct {
	service usecase.UserUsecase
}

func NewUserHandler(service usecase.UserUsecase) UserHandler {
	return &userHandler{
		service: service,
	}
}

func (uh *userHandler) GetMyProfile(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	userID, _ := claims["id"].(string)

	profile, err := uh.service.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	res := &dto.UserProfileDTO{
		ID:          profile.ID,
		Email:       profile.Email,
		Name:        profile.Name,
		MiddleName:  profile.MiddleName,
		Surname:     profile.Surname,
		PhoneNumber: profile.PhoneNumber,
		IsVerified:  profile.IsVerified,
	}

	return c.JSON(res)
}

func (uh *userHandler) EditMyProfile(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	userID, _ := claims["id"].(string)

	req := dto.EditProfileDTO{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := uh.service.EditProfile(c.Context(), userID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}
