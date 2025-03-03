package repository

import (
	"github.com/hewpao/hewpao-backend/dto"
)

type NotificationRepository interface {
	Notify(req *dto.NotificationDTO) error
}
