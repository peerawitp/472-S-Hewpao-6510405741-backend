package repository

import (
	"errors"

	"github.com/hewpao/hewpao-backend/domain"
)

type NotificationRepository interface {
	PrNotify(user *domain.User, producRequest *domain.ProductRequest) error
}

type NotificationRepositoryFactory interface {
	Register(provider string, repo NotificationRepository)
	GetRepository(provider string) (NotificationRepository, error)
}

type notificationRepositoryFactory struct {
	repos map[string]NotificationRepository
}

func NewNotificationRepositoryFactory() NotificationRepositoryFactory {
	return &notificationRepositoryFactory{repos: make(map[string]NotificationRepository)}
}

func (f *notificationRepositoryFactory) Register(provider string, repo NotificationRepository) {
	f.repos[provider] = repo
}

func (f *notificationRepositoryFactory) GetRepository(provider string) (NotificationRepository, error) {
	repo, exists := f.repos[provider]
	if !exists {
		return nil, errors.New("unsupported notification provider")
	}
	return repo, nil
}
