package gorm

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/repository"
	"gorm.io/gorm"
)

type ChatRepo struct {
	db *gorm.DB
}

func NewChatRepo(db *gorm.DB) repository.ChatRepository {
	return &ChatRepo{db: db}
}

func (c *ChatRepo) Create(chat *domain.Chat) error {
	err := c.db.Create(&chat)

	if err != nil {
		return err.Error
	}
	return nil
}

func (c *ChatRepo) GetByID(id uint) (*domain.Chat, error) {
	var chat domain.Chat
	err := c.db.First(&chat, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (c *ChatRepo) GetByName(name string) (*domain.Chat) {
	var chat domain.Chat
	err := c.db.First(&chat, "name = ?", name).Error
	if err != nil {
		return nil
	}
	return &chat
}
