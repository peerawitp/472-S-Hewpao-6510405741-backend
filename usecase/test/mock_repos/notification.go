package mock_repos

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/stretchr/testify/mock"
)

type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) PrNotify(user *domain.User, productRequest *domain.ProductRequest) error {
	args := m.Called(user, productRequest)
	return args.Error(0)
}
