package usecase

import (

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/repository"
)

type ChatUseCase interface {
	CreateChat(name string) error
	GetByID(id uint) (*domain.Chat, error)
	GetByName(name string) (*domain.Chat)
}

type ChatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) ChatUseCase {
	return &ChatService{repo: repo}
}

func (c *ChatService) CreateChat(name string) error {

	if c.repo.GetByName(name) != nil {
		return nil
	}
	chat := domain.Chat{Name: name}
	err := c.repo.Create(&chat)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatService) GetByID(id uint) (*domain.Chat, error) {
	return c.repo.GetByID(id)
}

func (c *ChatService) GetByName(name string) (*domain.Chat) {
	if c.repo.GetByName(name) == nil {
		return nil
	}
	
	return c.repo.GetByName(name)
}