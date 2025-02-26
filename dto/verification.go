package dto

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hewpao/hewpao-backend/types"
)

type GetUserVerificationDTO struct {
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	MiddleName  *string    `json:"middle_name"`
	Surname     string     `json:"surname"`
	PhoneNumber *string    `json:"phone_num"`
	IsVerified  bool       `json:"is_verified"`
	CardImage   *string    `json:"card_image"`
	Role        types.Role `json:"role"`
}

type VerifyWithKYCDTO struct {
	CardImage *fiber.FormFile `form:"card-image" validate:"required"`
}

type UpdateUserVerificationDTO struct {
	Isverified bool `json:"is_verified"`
}
