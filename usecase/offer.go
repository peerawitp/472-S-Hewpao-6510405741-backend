package usecase

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
)

type OfferUsecase interface {
	CreateOffer(req *dto.CreateOfferDTO, userID string) error
}

type offerService struct {
	offerRepo          repository.OfferRepository
	productRequestRepo repository.ProductRequestRepository
	userRepo           repository.UserRepository
	ctx                context.Context
}

func NewOfferService(
	offerRepo repository.OfferRepository,
	productRequestRepo repository.ProductRequestRepository,
	userRepo repository.UserRepository,
	ctx context.Context,
) OfferUsecase {
	return &offerService{
		offerRepo:          offerRepo,
		productRequestRepo: productRequestRepo,
		userRepo:           userRepo,
		ctx:                ctx,
	}
}

func (o *offerService) CreateOffer(req *dto.CreateOfferDTO, userID string) error {
	prID := int(req.ProductRequestID)
	productRequest, err := o.productRequestRepo.FindByID(prID)
	if err != nil {
		return err
	}

	user, err := o.userRepo.FindByID(o.ctx, userID)
	if err != nil {
		return err
	}

	offer := domain.Offer{
		ProductRequestID: &req.ProductRequestID,
		ProductRequest:   productRequest,
		OfferDate:        req.OfferDate,
		UserID:           userID,
		User:             user,
	}

	err = o.offerRepo.Create(&offer)
	if err != nil {
		return err
	}

	return nil
}
