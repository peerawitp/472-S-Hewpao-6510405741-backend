package mock_repos

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/stretchr/testify/mock"
)

type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) PrNotify(user *domain.User, productRequest *domain.ProductRequest) error {
	args := m.Called(user, productRequest)
	return args.Error(0)
}

type MockNotificationRepositoryFactory struct {
	mock.Mock
}

func (m *MockNotificationRepositoryFactory) Register(provider string, repo repository.NotificationRepository) {
	m.Called(provider, repo)
}

func (m *MockNotificationRepositoryFactory) GetRepository(provider string) (repository.NotificationRepository, error) {
	args := m.Called(provider)
	if repo, ok := args.Get(0).(repository.NotificationRepository); ok {
		return repo, args.Error(1)
	}
	return nil, args.Error(1)
}
