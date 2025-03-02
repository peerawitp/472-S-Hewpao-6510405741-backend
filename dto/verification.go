package dto

import (
	"github.com/gofiber/fiber/v2"
)

type VerifyWithKYCDTO struct {
	CardImage *fiber.FormFile `form:"card-image" validate:"required"`
}

type UpdateUserVerificationDTO struct {
	Isverified bool `json:"is_verified"`
}
