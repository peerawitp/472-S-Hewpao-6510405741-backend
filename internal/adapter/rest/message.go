package rest

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
)

type MessageHandler interface {
	CreateMessage(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx)(*domain.Message, error)
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

func (m *messageHandler)GetByID(c *fiber.Ctx) (*domain.Message, error) {
	id := c.Params("id")
	res, err := m.service.GetByID(id)
	if err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return res, nil
}

