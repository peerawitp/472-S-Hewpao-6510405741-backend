package repository

import (
	"errors"

	"github.com/hewpao/hewpao-backend/domain"
)

type NotificationRepository interface {
	PrNotify(user *domain.User, producRequest *domain.ProductRequest) error
}

type NotificationRepositoryFactory struct {
	repos map[string]NotificationRepository
}

func NewNotificationRepositoryFactory() NotificationRepositoryFactory {
	return NotificationRepositoryFactory{repos: make(map[string]NotificationRepository)}
}

func (f *NotificationRepositoryFactory) Register(provider string, repo NotificationRepository) {
	f.repos[provider] = repo
}

func (f *NotificationRepositoryFactory) GetRepository(provider string) (NotificationRepository, error) {
	repo, exists := f.repos[provider]
	if !exists {
		return nil, errors.New("unsupported notification provider")
	}
	return repo, nil
}
