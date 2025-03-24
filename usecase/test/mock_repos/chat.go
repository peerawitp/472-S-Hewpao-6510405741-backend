package mock_repos

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/stretchr/testify/mock"
)

type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) Create(chat *domain.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}

func (m *MockChatRepository) GetByID(id uint) (*domain.Chat, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatRepository) GetByName(name string) *domain.Chat {
	args := m.Called(name)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Chat)
	}
	return nil
}
