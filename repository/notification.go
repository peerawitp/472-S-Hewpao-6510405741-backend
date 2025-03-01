package repository

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
)

type NotificationRepository interface {
	Notify(toUser *domain.User, req *dto.NotificationDTO) error
}
