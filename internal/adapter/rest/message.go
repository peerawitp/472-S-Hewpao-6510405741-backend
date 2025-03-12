package rest

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
)

type MessageHandler interface {
	CreateMessage(c *fiber.Ctx) error
	GetByChatID(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
}

type messageHandler struct {
	service usecase.MessageService
}

func NewMessageHandler(service usecase.MessageService) MessageHandler {
	return &messageHandler{
		service: service,
	}
}

func (m *messageHandler) CreateMessage(c *fiber.Ctx) error {
	req := dto.CreateMessageRequestDTO{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Println(req.UserID)
	log.Println(req.ChatID)
	log.Println("*------------------------------------------------------------------------------")

	message, err := m.service.CreateMessage(req.UserID, req.ChatID, req.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(message)
}

func (m *messageHandler) GetByChatID(c *fiber.Ctx) error {
	chatID := c.Params("chat_id")

	messages, err := m.service.GetByChatID(chatID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(messages)
}

func (m *messageHandler) GetByID(c *fiber.Ctx) error {
	messageID := c.Params("id") // Make sure this matches your route param

	message, err := m.service.GetByID(messageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(message)
}
