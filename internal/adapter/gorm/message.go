package gorm

import (
	"gorm.io/gorm"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/hewpao/hewpao-backend/domain"
)

type MessageGormRepo struct {
	db *gorm.DB
}

func NewMessageGormRepo(db *gorm.DB) repository.MessageRepository {
	return &MessageGormRepo{db: db}
}

func (m *MessageGormRepo) Store(message *domain.Message) error {
	return m.db.Create(message).Error
	
}

func (m *MessageGormRepo) GetByID(id string) (*domain.Message, error) {
	var message domain.Message
	result := m.db.First(&message, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &message, nil
}

func (m *MessageGormRepo)GetByChatID(id string) ([]domain.Message, error) {
	var messages []domain.Message
	result := m.db.Where("chat_id = ?", id).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}