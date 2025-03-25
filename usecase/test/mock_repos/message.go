package mock_repos

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/stretchr/testify/mock"
)

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Store(message *domain.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockMessageRepository) GetByChatID(id string) ([]domain.Message, error) {
	args := m.Called(id)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageRepository) GetByID(id string) (*domain.Message, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Message), args.Error(1)
}
