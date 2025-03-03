package usecase

import (
	"context"
	"html/template"
	"strings"
	"time"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/hewpao/hewpao-backend/types"
	"gopkg.in/gomail.v2"
)

type NotificationUsecase interface {
	PrNotify(prod *domain.ProductRequest) error
}

type productRequestNotifier struct {
	repo      repository.NotificationRepository
	userRepo  repository.UserRepository
	ctx       context.Context
	message   *gomail.Message
	req       *domain.ProductRequest
	cfg       *config.Config
	offerRepo repository.OfferRepository
}

func NewNotificationUsecase(repo repository.NotificationRepository, userRepo repository.UserRepository, ctx context.Context, message *gomail.Message, cfg *config.Config, offerRepo repository.OfferRepository) NotificationUsecase {
	return &productRequestNotifier{
		repo:      repo,
		userRepo:  userRepo,
		offerRepo: offerRepo,
		ctx:       ctx,
		message:   message,
		cfg:       cfg,
	}
}

func sendToReciever(toUserID string, pn *productRequestNotifier, prod *domain.ProductRequest) error {
	user, err := pn.userRepo.FindByID(pn.ctx, toUserID)
	if err != nil {
		return err
	}

	var content strings.Builder

	err = notifyUpdate(prod, pn.cfg, &content)
	if err != nil {
		return err
	}

	notification := dto.NotificationDTO{
		ToEmail: user.Email,
		Subject: "[HEWPAO] Product Request Current Status Report",
		Content: content.String(),
	}

	err = pn.repo.Notify(&notification)
	if err != nil {
		return err
	}

	return nil
}

func (pn *productRequestNotifier) PrNotify(prod *domain.ProductRequest) error {
	var toUserID string

	offer := domain.Offer{}
	offer.ID = *prod.SelectedOfferID
	err := pn.offerRepo.GetByID(&offer)
	if err != nil {
		return err
	}

	switch prod.DeliveryStatus {
	case types.Pending, types.Purchased, types.Refunded:
		toUserID = *prod.UserID
		err := sendToReciever(toUserID, pn, prod)
		if err != nil {
			return err
		}
	case types.Opening, types.Cancel, types.Returned:
		toUserID = offer.UserID
		err = sendToReciever(toUserID, pn, prod)
		if err != nil {
			return err
		}
	case types.PickedUp, types.OutForDelivery, types.Delivered:
		buyerID := *prod.UserID
		travelerID := offer.UserID
		err = sendToReciever(buyerID, pn, prod)
		if err != nil {
			return err
		}
		err = sendToReciever(travelerID, pn, prod)
		if err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}

func notifyUpdate(productRequest *domain.ProductRequest, cfg *config.Config, content *strings.Builder) error {
	// Prepare the notification data
	data := dto.NotificationDataDTO{
		RecipientName: productRequest.User.Name,
		CompanyName:   "HEWPAO",
		ProductID:     productRequest.ID,
		ProductStatus: productRequest.DeliveryStatus,
		SupportEmail:  cfg.EmailUser,
		Year:          time.Now().Year(),
	}

	// Parse the template file
	tmpl, err := template.ParseFiles("./assets/emailTemplate.html")
	if err != nil {
		return err
	}

	// Clear any existing content in the builder before executing
	content.Reset()

	// Execute the template and build content
	err = tmpl.Execute(content, data)
	if err != nil {
		return err
	}

	return nil
}
