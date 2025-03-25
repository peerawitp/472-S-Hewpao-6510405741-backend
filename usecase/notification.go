package usecase

import (
	"context"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/repository"

	"github.com/hewpao/hewpao-backend/types"
	"gopkg.in/gomail.v2"
)

type NotificationUsecase interface {
	PrNotify(prod *domain.ProductRequest, provider string) error
}

type productRequestNotifier struct {
	notificationRepoFactory repository.NotificationRepositoryFactory
	userRepo                repository.UserRepository
	ctx                     context.Context
	message                 *gomail.Message
	cfg                     *config.Config
	offerRepo               repository.OfferRepository
}

func NewNotificationUsecase(notificationRepoFactory repository.NotificationRepositoryFactory, userRepo repository.UserRepository, ctx context.Context, message *gomail.Message, cfg *config.Config, offerRepo repository.OfferRepository) NotificationUsecase {
	return &productRequestNotifier{
		notificationRepoFactory: notificationRepoFactory,
		userRepo:                userRepo,
		offerRepo:               offerRepo,
		ctx:                     ctx,
		message:                 message,
		cfg:                     cfg,
	}
}

func prSend(toUserID string, pn *productRequestNotifier, prod *domain.ProductRequest, notificationRepo repository.NotificationRepository) error {
	user, err := pn.userRepo.FindByID(pn.ctx, toUserID)
	if err != nil {
		return err
	}

	err = notificationRepo.PrNotify(user, prod)
	if err != nil {
		return err
	}

	return nil
}

func (pn *productRequestNotifier) PrNotify(prod *domain.ProductRequest, provider string) error {
	notificationRepo, err := pn.notificationRepoFactory.GetRepository(provider)
	if err != nil {
		return err
	}

	var toUserID string

	offer := domain.Offer{}
	if prod.SelectedOfferID != nil {
		offer.ID = *prod.SelectedOfferID
		err = pn.offerRepo.GetByID(&offer)
		if err != nil {
			return err
		}
	}

	switch prod.DeliveryStatus {
	case types.Pending, types.Purchased, types.Refunded:
		toUserID = *prod.UserID
		err := prSend(toUserID, pn, prod, notificationRepo)
		if err != nil {
			return err
		}
	case types.Opening, types.Cancel, types.Returned:
		toUserID = offer.UserID
		if toUserID == "" {
			return nil
		}

		err = prSend(toUserID, pn, prod, notificationRepo)
		if err != nil {
			return err
		}
	case types.PickedUp, types.OutForDelivery, types.Delivered:
		buyerID := *prod.UserID

		err = prSend(buyerID, pn, prod, notificationRepo)
		if err != nil {
			return err
		}

		travelerID := offer.UserID
		if travelerID == "" {
			return nil
		}

		err = prSend(travelerID, pn, prod, notificationRepo)
		if err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}
