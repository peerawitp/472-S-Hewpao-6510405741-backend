package usecase

import (
	"context"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
)

type OfferUsecase interface {
	CreateOffer(req *dto.CreateOfferDTO, userID string) error
	GetOfferDetailByOfferID(offerID int, buyerUserID string) (*dto.GetOfferDetailDTO, error)
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

	if userID == *productRequest.UserID {
		return exception.ErrCouldNotSelfOffer
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

func (o *offerService) GetOfferDetailByOfferID(offerID int, buyerUserID string) (*dto.GetOfferDetailDTO, error) {
	offer, err := o.offerRepo.GetOfferDetailByOfferID(offerID)
	if err != nil {
		return nil, err
	}

	isUserOwnPR, err := o.productRequestRepo.IsOwnedByUser(int(offer.ProductRequest.ID), buyerUserID)
	if err != nil {
		return nil, err
	}

	if !isUserOwnPR {
		return nil, exception.ErrPermissionDenied
	}

	res := &dto.GetOfferDetailDTO{
		ID:               offer.ID,
		ProductRequestID: offer.ProductRequestID,
		ProductRequest:   offer.ProductRequest,
		UserID:           offer.UserID,
		User:             offer.User,
		OfferDate:        offer.OfferDate,
	}

	return res, nil
}
