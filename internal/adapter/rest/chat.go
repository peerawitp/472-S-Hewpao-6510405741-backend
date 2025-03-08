package rest

import (

	"github.com/gofiber/fiber/v2"
	"github.com/hewpao/hewpao-backend/usecase"
)

type ChatHandler interface {
	CreateChat(ctx *fiber.Ctx) error
}

type chatHandler struct {
	service usecase.ChatUseCase
}

func NewChatHandler(service usecase.ChatUseCase) ChatHandler {
	return &chatHandler{
		service: service,
	}
}

func (c *chatHandler) CreateChat(ctx *fiber.Ctx) error {

	var request struct {
		Name    string `json:"name"`
		OrderID uint `json:"order_id"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	err := c.service.CreateChat(request.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Chat created successfully",
	})
}
