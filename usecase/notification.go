package usecase

import (
	"context"

	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"gopkg.in/gomail.v2"
)

type NotificationUsecase interface {
	Notify(req *dto.NotificationDTO) error
}

type notificationService struct {
	repo     repository.NotificationRepository
	userRepo repository.UserRepository
	ctx      context.Context
	message  *gomail.Message
}

func NewNotificationUsecase(repo repository.NotificationRepository, userRepo repository.UserRepository, ctx context.Context, message *gomail.Message) NotificationUsecase {
	return &notificationService{
		repo:     repo,
		userRepo: userRepo,
		ctx:      ctx,
		message:  message,
	}
}

func (n *notificationService) Notify(req *dto.NotificationDTO) error {
	toUser, err := n.userRepo.FindByID(n.ctx, req.ToID)
	if err != nil {
		return err
	}

	err = n.repo.Notify(toUser, req)
	if err != nil {
		return err
	}

	return nil
}
