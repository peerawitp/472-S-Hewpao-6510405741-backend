package usecase

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/repository"
)

type MessageService struct {
	repo repository.MessageRepository
}

type MessageUsecase interface {
	CreateMessage(userID string, chatID string, text string) (*domain.Message, error)
	GetByChatID(id string) ([]domain.Message, error)
	GetByID(id string) (*domain.Message, error)
}

func NewMessageService(ms repository.MessageRepository) *MessageService {
	return &MessageService{repo: ms}
}

func (m *MessageService) CreateMessage(userID string, chatID uint, text string) (*domain.Message, error) {

	message := &domain.Message{
		UserID:  userID,
		ChatID:  chatID,
		Content: text,
	}

	err := m.repo.Store(message)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (m *MessageService) GetByChatID(id string) ([]domain.Message, error) {

	var messages []domain.Message
	messages, err := m.repo.GetByChatID(id)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *MessageService) GetByID(id string) (*domain.Message, error) {
	var message *domain.Message
	message, err := m.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return message, nil
}
