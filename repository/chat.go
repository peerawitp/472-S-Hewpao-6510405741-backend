package repository

import (
	"github.com/hewpao/hewpao-backend/domain"
)

type ChatRepository interface {
	Create(chat *domain.Chat) error
	GetByID(id uint) (*domain.Chat, error)
	GetByName(name string) (*domain.Chat)
}