package repository

import(
	"github.com/hewpao/hewpao-backend/domain"
)

type MessageRepository interface {
	Store(message *domain.Message) error
	GetByChatID(id string) ([]domain.Message, error)
	GetByID(id string) (*domain.Message, error)
}


